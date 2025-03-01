package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinReg(t *testing.T) {
	cases := []struct {
		CaseDescription string
		Path            string
		Fallback        string
		ExpectedSuccess bool
		ExpectedValue   string
		getWRKVOutput   *windowsRegistryValue
		Err             error
	}{
		{
			CaseDescription: "Error",
			Path:            "HKLLM\\Software\\Microsoft\\Windows NT\\CurrentVersion\\ProductName",
			Err:             errors.New("No match"),
			ExpectedSuccess: false,
		},
		{
			CaseDescription: "Value",
			Path:            "HKLM\\Software\\Microsoft\\Windows NT\\CurrentVersion\\InstallTime",
			getWRKVOutput:   &windowsRegistryValue{valueType: regString, str: "xbox"},
			ExpectedSuccess: true,
			ExpectedValue:   "xbox",
		},
		{
			CaseDescription: "Fallback value",
			Path:            "HKLM\\Software\\Microsoft\\Windows NT\\CurrentVersion\\InstallTime",
			Fallback:        "cortana",
			Err:             errors.New("No match"),
			ExpectedSuccess: true,
			ExpectedValue:   "cortana",
		},
		{
			CaseDescription: "Empty string value (no error) should display empty string even in presence of fallback",
			Path:            "HKLM\\Software\\Microsoft\\Windows NT\\CurrentVersion\\InstallTime",
			getWRKVOutput:   &windowsRegistryValue{valueType: regString, str: ""},
			Fallback:        "anaconda",
			ExpectedSuccess: true,
			ExpectedValue:   "",
		},
		{
			CaseDescription: "Empty string value (no error) should display empty string",
			Path:            "HKLM\\Software\\Microsoft\\Windows NT\\CurrentVersion\\InstallTime",
			getWRKVOutput:   &windowsRegistryValue{valueType: regString, str: ""},
			ExpectedSuccess: true,
			ExpectedValue:   "",
		},
		{
			CaseDescription: "DWORD value",
			Path:            "HKLM\\Software\\Microsoft\\Windows NT\\CurrentVersion\\InstallTime",
			getWRKVOutput:   &windowsRegistryValue{valueType: regDword, dword: 0xdeadbeef},
			ExpectedSuccess: true,
			ExpectedValue:   "0xDEADBEEF",
		},
		{
			CaseDescription: "QWORD value",
			Path:            "HKLM\\Software\\Microsoft\\Windows NT\\CurrentVersion\\InstallTime",
			getWRKVOutput:   &windowsRegistryValue{valueType: regQword, qword: 0x7eb199e57fa1afe1},
			ExpectedSuccess: true,
			ExpectedValue:   "0x7EB199E57FA1AFE1",
		},
	}

	for _, tc := range cases {
		env := new(MockedEnvironment)
		env.On("getRuntimeGOOS").Return(windowsPlatform)
		env.On("getWindowsRegistryKeyValue", tc.Path).Return(tc.getWRKVOutput, tc.Err)
		env.onTemplate()
		r := &winreg{
			env: env,
			props: properties{
				RegistryPath: tc.Path,
				Fallback:     tc.Fallback,
			},
		}

		assert.Equal(t, tc.ExpectedSuccess, r.enabled(), tc.CaseDescription)
		assert.Equal(t, tc.ExpectedValue, r.string(), tc.CaseDescription)
	}
}
