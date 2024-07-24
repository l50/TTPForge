/*
Copyright © 2024-present, Meta Platforms, Inc. and affiliates
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package blocks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func expectNoError(t *testing.T) expect.ConsoleOpt {
	return expect.WithExpectObserver(
		func(matchers []expect.Matcher, buf string, err error) {
			if err == nil {
				return
			}
			if len(matchers) == 0 {
				t.Fatalf("Error occurred while matching %q: %s\n%s", buf, err, string(debug.Stack()))
			} else {
				var criteria []string
				for _, matcher := range matchers {
					if crit, ok := matcher.Criteria().([]string); ok {
						criteria = append(criteria, crit...)
					} else {
						criteria = append(criteria, "unknown criteria")
					}
				}
				t.Fatalf("Failed to find [%s] in %q: %s\n%s", strings.Join(criteria, ", "), buf, err, string(debug.Stack()))
			}
		},
	)
}

func createTestScript(t *testing.T, scriptContent string) (string, string) {
	tempDir, err := os.MkdirTemp("", "python-script-test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	scriptPath := filepath.Join(tempDir, "interactive.py")
	err = os.WriteFile(scriptPath, []byte(scriptContent), 0644)
	require.NoError(t, err)
	return scriptPath, tempDir
}

func sendNoError(t *testing.T) expect.ConsoleOpt {
	return expect.WithSendObserver(
		func(msg string, n int, err error) {
			if err != nil {
				t.Fatalf("Failed to send %q: %s\n%s", msg, err, string(debug.Stack()))
			}
			if len(msg) != n {
				t.Fatalf("Only sent %d of %d bytes for %q\n%s", n, len(msg), msg, string(debug.Stack()))
			}
		},
	)
}

func NewTestTTPExecutionContext(workDir string) TTPExecutionContext {
	return TTPExecutionContext{
		WorkDir: workDir,
	}
}

// func TestExpectStep(t *testing.T) {
// 	t.Parallel()
// 	testCases := []struct {
// 		name               string
// 		script             string
// 		content            string
// 		wantUnmarshalError bool
// 		wantValidateError  bool
// 		wantExecuteError   bool
// 		expectedErrTxt     string
// 	}{
// 		{
// 			name: "Test Unmarshal Expect Valid",
// 			script: `
// print("Enter your name:")
// name = input()
// print(f"Hello, {name}!")
// print("Enter your age:")
// age = input()
// print(f"You are {age} years old.")
// `,
// 			content: `
// steps:
//   - name: run_expect_script
//     description: "Run an expect script to interact with the command."
//     expect:
//       inline: |
//         python3 interactive.py
//       responses:
//         - prompt: "Enter your name:"
//           response: "John"
//         - prompt: "Enter your age:"
//           response: "30"
// `,
// 		},
// 		{
// 			name: "Test Unmarshal Expect No Inline",
// 			script: `
// print("Enter your name:")
// name = input()
// print(f"Hello, {name}!")
// print("Enter your age:")
// age = input()
// print(f"You are {age} years old.")
// `,
// 			content: `
// steps:
//   - name: run_expect_script
//     description: "Run an expect script to interact with the command."
//     expect:
//       responses:
//         - prompt: "Enter your name:"
//           response: "John"
//         - prompt: "Enter your age:"
//           response: "30"
// `,
// 			wantValidateError: true,
// 			expectedErrTxt:    "inline must be provided",
// 		},
// 		{
// 			name: "Test ExpectStep Execute With Output",
// 			script: `
// print("Enter your name:")
// name = input()
// print(f"Hello, {name}!")
// print("Enter your age:")
// age = input()
// print(f"You are {age} years old.")
// `,
// 			content: `
// steps:
//   - name: run_expect_script
//     description: "Run an expect script to interact with the command."
//     expect:
//       inline: |
//         python3 interactive.py
//       responses:
//         - prompt: "Enter your name:"
//           response: "John"
//         - prompt: "Enter your age:"
//           response: "30"
// `,
// 		},
// 		{
// 			name: "Test ExpectStep with Chdir",
// 			script: `
// import os
// print("Current directory:", os.getcwd())
// print("Enter a number:")
// number = input()
// print(f"You input {number}.")
// `,
// 			content: `
// steps:
//   - name: run_expect_script
//     description: "Run an expect script to interact with the command."
//     expect:
//       chdir: "/tmp"
//       inline: |
//         python3 interactive.py
//       responses:
//         - prompt: "Enter a number:"
//           response: "30"
// `,
// 		},
// 		{
// 			name: "Test ExpectStep with CleanupStep",
// 			script: `
// print("Enter your name:")
// name = input()
// print(f"Hello, {name}!")
// print("Enter your age:")
// age = input()
// print(f"You are {age} years old.")
// `,
// 			content: `
// steps:
//   - name: run_expect_script
//     description: "Run an expect script to interact with the command."
//     expect:
//       chdir: "/tmp"
//       inline: |
//         python3 interactive.py
//       responses:
//         - prompt: "Enter your name:"
//           response: "John"
//         - prompt: "Enter your age:"
//           response: "30"
//       cleanup: |
//         pwd
//         cat interactive.py
//         rm interactive.py
// `,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			scriptPath, tempDir := createTestScript(t, tc.script)

// 			var steps struct {
// 				Steps []struct {
// 					Name        string      `yaml:"name"`
// 					Description string      `yaml:"description"`
// 					Expect      *ExpectStep `yaml:"expect"`
// 				} `yaml:"steps"`
// 			}

// 			err := yaml.Unmarshal([]byte(tc.content), &steps)
// 			if tc.wantUnmarshalError {
// 				assert.Error(t, err)
// 				return
// 			}
// 			require.NoError(t, err)

// 			if len(steps.Steps) == 0 || steps.Steps[0].Expect == nil {
// 				assert.Fail(t, "Failed to unmarshal test case content")
// 				return
// 			}

// 			expectStep := steps.Steps[0].Expect

// 			err = expectStep.Validate(NewTestTTPExecutionContext(tempDir))
// 			if tc.wantValidateError {
// 				assert.Equal(t, tc.expectedErrTxt, err.Error())
// 				return
// 			}
// 			require.NoError(t, err)

// 			execCtx := NewTestTTPExecutionContext(tempDir)
// 			console, err := expect.NewConsole(expect.WithStdout(os.Stdout), expect.WithStdin(os.Stdin))
// 			require.NoError(t, err)
// 			defer console.Close()

// 			cmd := exec.Command("sh", "-c", "python3 "+scriptPath)
// 			cmd.Stdin = console.Tty()
// 			cmd.Stdout = console.Tty()
// 			cmd.Stderr = console.Tty()

// 			if expectStep.Chdir != "" {
// 				cmd.Dir = expectStep.Chdir
// 			}

// 			err = cmd.Start()
// 			require.NoError(t, err)

// 			done := make(chan struct{})

// 			go func() {
// 				defer close(done)
// 				for _, response := range expectStep.Expect.Responses {
// 					re := regexp.MustCompile(response.Prompt)
// 					_, err := console.Expect(expect.Regexp(re), expect.WithTimeout(60*time.Second))
// 					if err != nil {
// 						t.Errorf("failed to expect %q: %v", re, err)
// 						return
// 					}
// 					_, err = console.SendLine(response.Response)
// 					if err != nil {
// 						t.Errorf("failed to send response: %v", err)
// 						return
// 					}
// 				}
// 				console.Tty().Close()
// 			}()

// 			_, err = expectStep.Execute(execCtx)
// 			if tc.wantExecuteError {
// 				assert.Error(t, err)
// 				assert.Contains(t, err.Error(), tc.expectedErrTxt)
// 				return
// 			}
// 			require.NoError(t, err)

// 			<-done

// 			output, err := console.ExpectEOF()
// 			require.NoError(t, err)

// 			normalizedOutput := strings.ReplaceAll(output, "\r\n", "\n")

// 			if tc.name == "Test ExpectStep with Chdir" {
// 				assert.Contains(t, normalizedOutput, "Current directory: ")
// 				assert.Contains(t, normalizedOutput, "You input 30.\n")
// 			}

// 			if tc.name == "Test ExpectStep with CleanupStep" {
// 				assert.Contains(t, normalizedOutput, "Hello, John!\n")
// 				assert.Contains(t, normalizedOutput, "You are 30 years old.\n")

// 				result, err := expectStep.Cleanup(execCtx)
// 				require.NoError(t, err)
// 				assert.NotNil(t, result)
// 			}
// 		})
// 	}
// }

// func TestYamlUnmarshal(t *testing.T) {
// 	yamlConfig := `
// steps:
//   - name: run_expect_script
//     expect:
//       inline: |
//         echo "Hello"
//       responses:
//         - prompt: "Enter your name: "
//           response: "John"
// `

// 	var config struct {
// 		Steps []struct {
// 			Name   string      `yaml:"name"`
// 			Expect *ExpectStep `yaml:"expect"`
// 		} `yaml:"steps"`
// 	}

// 	err := yaml.Unmarshal([]byte(yamlConfig), &config)
// 	require.NoError(t, err)
// 	require.Len(t, config.Steps, 1)
// 	require.NotNil(t, config.Steps[0].Expect)

// 	t.Logf("Unmarshaled config: %+v", config)
// }

func TestExpectStepSSH(t *testing.T) {
	t.Parallel()
	password := "Password123!"

	// Bash script to simulate SSH login
	script := fmt.Sprintf(`#!/bin/bash
echo -n "bobbo@k8s6's password: "
read -r pwd
if [ "$pwd" == "%s" ]; then
    while true; do
        echo -n "bobbo@k8s6:~$ "
        read -r cmd
        if [ "$cmd" == "whoami" ]; then
            echo "bobbo"
            exit 0
        else
            echo "Unknown command"
        fi
    done
else
    echo "Authentication failed."
    exit 1
fi
`, password)

	scriptPath, tempDir := createTestScript(t, script)
	defer os.RemoveAll(tempDir)

	// Make the script executable
	err := os.Chmod(scriptPath, 0755)
	require.NoError(t, err)

	yamlConfig := fmt.Sprintf(`
steps:
  - name: run_expect_script
    expect:
      inline: |
        sh %s
      responses:
        - prompt: "bobbo@k8s6's password: "
          response: "%s"
        - prompt: "bobbo@k8s6:~$ "
          response: "whoami"
`, scriptPath, password)

	var config struct {
		Steps []struct {
			Name   string      `yaml:"name"`
			Expect *ExpectStep `yaml:"expect"`
		} `yaml:"steps"`
	}

	err = yaml.Unmarshal([]byte(yamlConfig), &config)
	require.NoError(t, err)
	require.Len(t, config.Steps, 1)
	require.NotNil(t, config.Steps[0].Expect)

	expectStep := config.Steps[0].Expect

	// Add debug information
	t.Logf("ExpectStep: %+v", expectStep)

	execCtx := NewTestTTPExecutionContext(tempDir)

	err = expectStep.Validate(execCtx)
	require.NoError(t, err)

	console, err := expect.NewConsole(expectNoError(t), sendNoError(t), expect.WithStdout(os.Stdout), expect.WithStdin(os.Stdin))
	require.NoError(t, err)
	defer console.Close()

	cmd := exec.Command("sh", scriptPath)
	cmd.Stdin = console.Tty()
	cmd.Stdout = console.Tty()
	cmd.Stderr = console.Tty()

	err = cmd.Start()
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() {
		for _, response := range expectStep.Expect.Responses {
			t.Logf("Expecting prompt: %s", response.Prompt)
			re := regexp.MustCompile(response.Prompt)
			matchedText, err := console.Expect(expect.Regexp(re))
			if err != nil {
				done <- fmt.Errorf("failed to expect %q: %w", re, err)
				return
			}
			t.Logf("Matched text: %s", matchedText)
			t.Logf("Sending response: %s", response.Response)
			_, err = console.SendLine(response.Response)
			if err != nil {
				done <- fmt.Errorf("failed to send response: %w", err)
				return
			}
		}
		console.Tty().Close() // Ensure TTY is closed to signal EOF
		done <- nil
	}()

	timeout := 30 * time.Second

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("error in expect: %v", err)
		}
	case <-time.After(timeout):
		t.Fatalf("timeout waiting for expect")
	}

	err = cmd.Wait()
	require.NoError(t, err)

	output, err := console.ExpectEOF()
	require.NoError(t, err)

	normalizedOutput := strings.ReplaceAll(output, "\r\n", "\n")
	expectedSubstring := "bobbo\n"

	t.Logf("Full output: %s", normalizedOutput)
	assert.Contains(t, normalizedOutput, expectedSubstring)
}

func TestExpectStepy(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name               string
		script             string
		content            string
		wantUnmarshalError bool
		wantValidateError  bool
		wantExecuteError   bool
		expectedErrTxt     string
	}{
		{
			name: "Test Unmarshal Expect Valid",
			script: `
print("Enter your name:")
name = input()
print(f"Hello, {name}!")
print("Enter your age:")
age = input()
print(f"You are {age} years old.")
`,
			content: `
steps:
  - name: run_expect_script
    description: "Run an expect script to interact with the command."
    expect:
      inline: |
        python3 interactive.py
      responses:
        - prompt: "Enter your name:"
          response: "John"
        - prompt: "Enter your age:"
          response: "30"
`,
		},
		{
			name: "Test Unmarshal Expect No Inline",
			script: `
print("Enter your name:")
name = input()
print(f"Hello, {name}!")
print("Enter your age:")
age = input()
print(f"You are {age} years old.")
`,
			content: `
steps:
  - name: run_expect_script
    description: "Run an expect script to interact with the command."
    expect:
      responses:
        - prompt: "Enter your name:"
          response: "John"
        - prompt: "Enter your age:"
          response: "30"
`,
			wantValidateError: true,
			expectedErrTxt:    "expect block must be provided",
		},
		{
			name: "Test ExpectStep Execute With Output",
			script: `
print("Enter your name:")
name = input()
print(f"Hello, {name}!")
print("Enter your age:")
age = input()
print(f"You are {age} years old.")
`,
			content: `
steps:
  - name: run_expect_script
    description: "Run an expect script to interact with the command."
    expect:
      chdir: "/tmp"
      inline: |
        python3 interactive.py
      responses:
        - prompt: "Enter your name:"
          response: "John"
        - prompt: "Enter your age:"
          response: "30"
      cleanup: |
        pwd
        cat interactive.py
        rm interactive.py
`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			scriptPath, tempDir := createTestScript(t, tc.script)

			var steps struct {
				Steps []struct {
					Name        string      `yaml:"name"`
					Description string      `yaml:"description"`
					Expect      *ExpectStep `yaml:"expect"`
				} `yaml:"steps"`
			}

			err := yaml.Unmarshal([]byte(tc.content), &steps)
			if tc.wantUnmarshalError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			if len(steps.Steps) == 0 || steps.Steps[0].Expect == nil {
				assert.Fail(t, "Failed to unmarshal test case content")
				return
			}

			expectStep := steps.Steps[0].Expect

			err = expectStep.Validate(NewTestTTPExecutionContext(tempDir))
			if tc.wantValidateError {
				assert.Equal(t, tc.expectedErrTxt, err.Error())
				return
			}
			require.NoError(t, err)

			execCtx := NewTestTTPExecutionContext(tempDir)
			console, err := expect.NewConsole(expect.WithStdout(os.Stdout), expect.WithStdin(os.Stdin))
			require.NoError(t, err)
			defer console.Close()

			cmd := exec.Command("sh", "-c", "python3 "+scriptPath)
			cmd.Stdin = console.Tty()
			cmd.Stdout = console.Tty()
			cmd.Stderr = console.Tty()

			if expectStep.Chdir != "" {
				cmd.Dir = expectStep.Chdir
			}

			err = cmd.Start()
			require.NoError(t, err)

			done := make(chan struct{})

			go func() {
				defer close(done)
				for _, response := range expectStep.Expect.Responses {
					re := regexp.MustCompile(response.Prompt)
					_, err := console.Expect(expect.Regexp(re), expect.WithTimeout(60*time.Second))
					if err != nil {
						t.Errorf("failed to expect %q: %v", re, err)
						return
					}
					_, err = console.SendLine(response.Response)
					if err != nil {
						t.Errorf("failed to send response: %v", err)
						return
					}
				}
				console.Tty().Close()
			}()

			_, err = expectStep.Execute(execCtx)
			if tc.wantExecuteError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrTxt)
				return
			}
			require.NoError(t, err)

			<-done

			output, err := console.ExpectEOF()
			require.NoError(t, err)

			normalizedOutput := strings.ReplaceAll(output, "\r\n", "\n")

			if tc.name == "Test ExpectStep Execute With Output" {
				assert.Contains(t, normalizedOutput, "Hello, John!\n")
				assert.Contains(t, normalizedOutput, "You are 30 years old.\n")

				result, err := expectStep.Cleanup(execCtx)
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestExpectStepr(t *testing.T) {
	// Ensure sshpass is in the PATH
	path := os.Getenv("PATH")
	sshpassPath := "/opt/homebrew/bin"
	os.Setenv("PATH", sshpassPath+":"+path)

	tests := []struct {
		name    string
		step    *ExpectStep
		wantErr bool
	}{
		{
			name: "Test_Unmarshal_Expect_Valid",
			step: &ExpectStep{
				Expect: &ExpectSpec{
					Inline: "echo 'hello world'",
					Responses: []Response{
						{Prompt: "hello", Response: "world"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test_Unmarshal_Expect_No_Inline",
			step: &ExpectStep{
				Expect: &ExpectSpec{
					Responses: []Response{
						{Prompt: "hello", Response: "world"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test_ExpectStep_Execute_With_Output",
			step: &ExpectStep{
				Executor: "bash",
				Expect: &ExpectSpec{
					Inline: `
					sshpass -p Password123! ssh bobbo@k8s6`,
					Responses: []Response{
						{Prompt: "Welcome to Ubuntu", Response: "whoami"},
						{Prompt: "bobbo", Response: "exit"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCtx := TTPExecutionContext{WorkDir: "."}
			fmt.Println("Executing command:", tt.step.Expect.Inline)
			_, err := tt.step.Execute(execCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("Command execution complete")
		})
	}
}
