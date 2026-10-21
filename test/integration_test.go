package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/gtkit/goerr"
)

// 外部包引用路径下的冒烟测试，确保导出 API 与文档一致。
func TestExportedAPI(t *testing.T) {
	item := goerr.New(errors.New("db down"), goerr.StatusMysqlServer(), "服务繁忙")
	if item.Message() == item.Error() {
		t.Error("Message() should differ from Error() when cause is present")
	}
	if item.Code() != goerr.ErrMysqlServer {
		t.Errorf("Code = %v", item.Code())
	}
	_, ok := goerr.AsItem(item)
	if !ok {
		t.Fatal("AsItem")
	}
}

// 仅用于调用方互操作验证，不是 goerr 的日志协议。
type diagnosticLogRecord struct {
	Operation  string     `json:"operation"`
	RequestID  string     `json:"request_id"`
	Code       goerr.Code `json:"code"`
	Error      string     `json:"error"`
	CauseCount int        `json:"cause_count"`
}

// 两个出口绑定各自的目标，返回同一种普通函数；不向 goerr 注册任何内容。
func textDiagnosticOutput(w io.Writer, operation, requestID string) func(goerr.Diagnostic) error {
	return func(diagnostic goerr.Diagnostic) error {
		_, err := fmt.Fprintf(w, "operation=%s request_id=%s code=%d causes=%d error=%s\n",
			operation, requestID, diagnostic.Code, len(diagnostic.Causes), goerr.RedactForLog(diagnostic.Error))
		return err
	}
}

func jsonDiagnosticOutput(w io.Writer, operation, requestID string) func(goerr.Diagnostic) error {
	return func(diagnostic goerr.Diagnostic) error {
		return json.NewEncoder(w).Encode(diagnosticLogRecord{
			Operation: operation, RequestID: requestID, Code: diagnostic.Code,
			Error: goerr.RedactForLog(diagnostic.Error), CauseCount: len(diagnostic.Causes),
		})
	}
}

func TestDiagnosticOutlets(t *testing.T) {
	t.Parallel()
	// 联系方式为合成夹具；出口需脱敏，原错误及原快照保留原文。
	cause := errors.New("call 13900001111 user@example.com")
	input := fmt.Errorf("load profile: %w", goerr.New(cause, goerr.StatusParams(), "invalid input"))
	diagnostic, err := goerr.Inspect(input)
	if err != nil {
		t.Fatal(err)
	}
	wantSnapshot, err := goerr.Inspect(input)
	if err != nil {
		t.Fatal(err)
	}
	wantBusiness := goerr.Convert(input)
	wantText := "operation=load_profile request_id=example-request code=10010002 causes=3 error=load profile: invalid input: call [phone] [email]\n"
	wantRecord := diagnosticLogRecord{
		Operation: "load_profile", RequestID: "example-request", Code: goerr.ErrParams,
		Error: "load profile: invalid input: call [phone] [email]", CauseCount: 3,
	}
	var textTarget, jsonTarget bytes.Buffer
	textOutput := textDiagnosticOutput(&textTarget, "load_profile", "example-request")
	jsonOutput := jsonDiagnosticOutput(&jsonTarget, "load_profile", "example-request")
	output := textOutput
	if err := output(diagnostic); err != nil {
		t.Fatal(err)
	}
	if textTarget.String() != wantText || jsonTarget.Len() != 0 {
		t.Fatal("selecting text must write only its target with explicit redaction")
	}
	output = jsonOutput
	if err := output(diagnostic); err != nil {
		t.Fatal(err)
	}
	var record diagnosticLogRecord
	if err := json.Unmarshal(jsonTarget.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record != wantRecord || textTarget.String() != wantText {
		t.Fatal("switching to JSON must preserve fields without writing the text target again")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(jsonTarget.Bytes(), &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 5 {
		t.Fatal("the caller projection must not serialize the complete snapshot or runtime metadata")
	}
	for _, emit := range []func(goerr.Diagnostic) error{textOutput, jsonOutput} {
		if err := emit(diagnostic); err != nil {
			t.Fatal(err)
		}
	}
	if textTarget.String() != strings.Repeat(wantText, 2) {
		t.Fatal("both ordinary functions must coexist without registration")
	}
	decoder := json.NewDecoder(bytes.NewReader(jsonTarget.Bytes()))
	for range 2 {
		if err := decoder.Decode(&record); err != nil || record != wantRecord {
			t.Fatal("each JSON event must contain the selected, redacted fields", err)
		}
	}
	if err := decoder.Decode(&record); !errors.Is(err, io.EOF) {
		t.Fatal("an outlet must not create extra events", err)
	}
	if !reflect.DeepEqual(diagnostic, wantSnapshot) || goerr.Convert(input) != wantBusiness ||
		!errors.Is(input, cause) || input.Error() != wantSnapshot.Error {
		t.Fatal("read-only outlets must not change the snapshot or original business error")
	}
}

type diagnosticFailureWriter struct {
	limit   int
	err     error
	calls   int
	written []byte
}

func (w *diagnosticFailureWriter) Write(data []byte) (int, error) {
	w.calls++
	n := min(w.limit, len(data))
	w.written = append(w.written, data[:n]...)
	return n, w.err
}

func TestDiagnosticOutletFailures(t *testing.T) {
	t.Parallel()
	for _, outlet := range []struct {
		name string
		bind func(io.Writer, string, string) func(goerr.Diagnostic) error
	}{
		{name: "text", bind: textDiagnosticOutput},
		{name: "json", bind: jsonDiagnosticOutput},
	} {
		for _, failure := range []struct {
			name  string
			limit int
		}{
			{name: "before write"},
			{name: "after partial write", limit: 7},
		} {
			t.Run(outlet.name+"/"+failure.name, func(t *testing.T) {
				t.Parallel()
				cause := errors.New("storage unavailable")
				input := fmt.Errorf("load: %w", goerr.New(cause, goerr.StatusNotFound(), "not found"))
				before := goerr.Convert(input)
				wantCode, wantHTTP, wantMessage, wantText := before.Code(), before.HTTPCode(), before.Message(), input.Error()
				diagnostic, err := goerr.Inspect(input)
				if err != nil {
					t.Fatal(err)
				}
				writer := &diagnosticFailureWriter{limit: failure.limit, err: io.ErrClosedPipe}
				output := outlet.bind(writer, "load_profile", "example-request")
				if err := output(diagnostic); !errors.Is(err, io.ErrClosedPipe) {
					t.Fatal("the caller must receive the output failure separately", err)
				}
				if writer.calls != 1 || len(writer.written) != failure.limit {
					t.Fatal("partial output must remain observable without automatic retry or rollback")
				}
				after := goerr.Convert(input)
				if after != before || after.Code() != wantCode || after.HTTPCode() != wantHTTP || after.Message() != wantMessage ||
					input.Error() != wantText || !errors.Is(input, cause) || errors.Is(input, io.ErrClosedPipe) ||
					!goerr.HasCode(input, goerr.ErrNotFound) {
					t.Fatal("output failure must not replace the business error or change response fields")
				}
			})
		}
	}
}

func TestDiagnosticFailureSkipsOutputs(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		input   error
		wantErr bool
	}{
		{name: "no business error"},
		{name: "inspection failure", input: errors.New(strings.Repeat("private diagnostic text ", 1000)), wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calls := 0
			output := func(goerr.Diagnostic) error { calls++; return nil }
			// 这是调用方的诊断操作，其 error 与 tt.input 分开处理。
			diagnosticErr := func() error {
				diagnostic, err := goerr.Inspect(tt.input)
				if err != nil {
					return err
				}
				if !diagnostic.HasError {
					return nil
				}
				return output(diagnostic)
			}()
			if (diagnosticErr != nil) != tt.wantErr || calls != 0 {
				t.Fatal("failed inspection or nil input must not call an outlet or fall back to raw text", diagnosticErr)
			}
			if tt.input != nil && (goerr.Convert(tt.input).Code() != goerr.ErrInternalServer ||
				!strings.HasPrefix(tt.input.Error(), "private diagnostic text ")) {
				t.Fatal("diagnostic failure must leave the original business error available")
			}
		})
	}
}

func TestDiagnosticConcurrentOutlets(t *testing.T) {
	t.Parallel()
	leftCause, rightCause := errors.New("call 13900001111"), errors.New("user@example.com")
	left := goerr.New(leftCause, goerr.StatusParams(), "invalid input")
	right := goerr.New(rightCause, goerr.StatusNotFound(), "not found")
	input := fmt.Errorf("request failed: %w", errors.Join(left, right))
	diagnostic, err := goerr.Inspect(input)
	if err != nil {
		t.Fatal(err)
	}
	before, err := goerr.Inspect(input)
	if err != nil {
		t.Fatal(err)
	}
	if diagnostic.Code != goerr.ErrParams || len(diagnostic.Causes) != 6 ||
		diagnostic.Causes[4].Code != goerr.ErrNotFound || goerr.HasCode(input, goerr.ErrNotFound) {
		t.Fatal("the primary status must retain the first Item while the other branch remains diagnostic data")
	}
	wantError := goerr.RedactForLog(input.Error())
	var workers sync.WaitGroup
	for range 16 {
		workers.Go(func() {
			// 同一快照只读共享；每个 worker 独占自己的 writer。
			var textTarget, jsonTarget bytes.Buffer
			for _, output := range []func(goerr.Diagnostic) error{
				textDiagnosticOutput(&textTarget, "load_profile", "example-request"),
				jsonDiagnosticOutput(&jsonTarget, "load_profile", "example-request"),
			} {
				if err := output(diagnostic); err != nil {
					t.Error(err)
					return
				}
			}
			var record diagnosticLogRecord
			if err := json.Unmarshal(jsonTarget.Bytes(), &record); err != nil {
				t.Error(err)
				return
			}
			if record.Code != goerr.ErrParams || record.CauseCount != 6 || record.Error != wantError ||
				!strings.Contains(textTarget.String(), wantError) || strings.Contains(textTarget.String(), "13900001111") {
				t.Error("concurrent outlets must preserve the selected fields and redact the original text")
			}
		})
	}
	workers.Wait()
	if !reflect.DeepEqual(diagnostic, before) || goerr.Convert(input) != left ||
		input.Error() != before.Error || !errors.Is(input, leftCause) || !errors.Is(input, rightCause) {
		t.Fatal("parallel read-only outputs must preserve the shared snapshot and original business error")
	}
}
