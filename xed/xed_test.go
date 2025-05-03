package xed_test

import (
	"github.com/ddkwork/ddk/xed"
	"testing"
)

func TestDecode32(t *testing.T) {
	data := []byte{
		// 779378E8 | EB 07                    | jmp ntdll.779378F1                      |
		// 779378EA | 33C0                     | xor eax,eax                             |
		// 779378EC | 40                       | inc eax                                 |
		// 779378ED | C3                       | ret                                     |
		// 779378EE | 8B65 E8                  | mov esp,dword ptr ss:[ebp-18]           | [dword ptr ss:[ebp-18]]:"6~稭"
		// 779378F1 | C745 FC FEFFFFFF         | mov dword ptr ss:[ebp-4],FFFFFFFE       |
		// 779378F8 | 8B4D F0                  | mov ecx,dword ptr ss:[ebp-10]           |
		// 779378FB | 64:890D 00000000         | mov dword ptr fs:[0],ecx                |
		// 77937902 | 59                       | pop ecx                                 |
		// 77937903 | 5F                       | pop edi                                 |
		// 77937904 | 5E                       | pop esi                                 |
		// 77937905 | 5B                       | pop ebx                                 |
		// 77937906 | C9                       | leave                                   |
		// 77937907 | C3                       | ret                                     |

		0x8B, 0x4D, 0xF0, 0x64, 0x89, 0x0D, 0x00, 0x00, 0x00, 0x00,
	}
	x := xed.New(data).Decode32()
	println(x.IntelSyntaxAsm.String())
}

func TestDecode64(t *testing.T) {
	// Decode64(tt.args.data)
}
