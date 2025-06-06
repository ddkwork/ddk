package winver

import (
	"github.com/ddkwork/ddk/xed"
	"github.com/hashicorp/go-version"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

// WindowVersion ps --> $PSVersionTable
func WindowVersion() (v string) {
	file := xed.ParserPe("C:\\Windows\\System32\\ntoskrnl.exe")
	resources := mylog.Check2(file.ParseVersionResources())

	for k, value := range resources {
		if k == "ProductVersion" {
			version := mylog.Check2(version.NewVersion(value))

			s := mylog.Check2(osVersion(value))

			s += " v"
			s += version.String()
			return s
			// mylog.Info(k, value)
		}
	}
	return
}

/*
Windows NT 4
Windows 95
Windows 98
Windows Me
Windows 2000
Windows XP
Windows XP 64
Windows Server 2003
Windows Server 2003 R2
Windwos Vista
Windows Server 2008
Windwos 7
Windows Server 2008 R2
Windwos 8
Windows Server 2012
Windows 8.1
Windows Server 2012 R2
Windows 10
Windows Server 2016
Windows Server 2019
Windows 11
Windows 11 +
*/
func osVersion(v string) (string, error) {
	version := stream.NewVersion(v)

	// 以下代码获取windows8.1及以上的操作系统
	majorVersion, minorVersion, buildNumber := version.Major, version.Minor, version.Patch
	// fmt.Printf("majorVersion:%d ,minorVersion:%d ,buildNumber:%d \n", majorVersion, minorVersion, buildNumber)

	o := mylog.Check2(GetVersionExW())

	// fmt.Println("GetVersionExW : ", err)

	/*
	   Windows 8.1
	   Windows 10
	   Windows Server 2016
	   Windows Server 2019
	   Windows 11
	   Windows 11 +
	*/
	// GetVersionExW : win8.1+ :> 6.2.9200 ()
	if majorVersion > 6 || (majorVersion == 6 && minorVersion >= 3) { // win8plus

		if majorVersion == 6 && minorVersion >= 3 {
			// Win8.1       : 6.3.9600
			// Windows Server 2012
			if o.ProductType == byte(VER_NT_WORKSTATION) {
				return "Windows 8.1", nil
			} else {
				// fmt.Println("o.ProductType :", o.ProductType)
				return "Windows Server 2012 R2", nil
			}
		} else if majorVersion == 10 && minorVersion == 0 {
			// Win 10       :10.0.19042
			// WinSer 2019  :10.0.17763
			// Win 11       :10.0.22000
			if o.ProductType == byte(VER_NT_WORKSTATION) {
				if buildNumber >= 22000 {
					return "Windows 11", nil
				} else { // if buildNumber >= 18363 {
					//  18363 : Win10 专业版
					//  19041 : win10 家庭中文版
					//  19042 : win10 家庭中文版,教育版
					//  19043 : win10 专业版
					return "Windows 10", nil
				}
			} else {
				if buildNumber >= 17763 {
					return "Windows Server 2019", nil
				} else if buildNumber >= 14393 {
					return "Windows Server 2016", nil
				}
			}
		} else {
			return "Windows 11 +", nil
		}
	}

	// 以下代码获取windows 8.1以下的系统版本
	/*
	   Windows NT 4
	   Windows 95
	   Windows 98
	   Windows Me
	   Windows 2000
	   Windows XP
	   Windows XP 64
	   Windows Server 2003
	   Windows Server 2003 R2
	   Windwos Vista
	   Windows Server 2008
	   Windwos 7
	   Windows Server 2008 R2
	   Windwos 8
	   Windows Server 2012
	*/

	// s, err := GetNativeSystemInfo()
	// if err != nil {
	//	return "", nil
	// }
	// u := s.GetDummyStructName()
	switch o.MajorVersion {
	case 4:
		{
			switch o.MinorVersion {
			case 0:
				{
					if int(o.PlatformId) == VER_PLATFORM_WIN32_NT {
						return "Windows NT 4", nil
					} else if int(o.PlatformId) == VER_PLATFORM_WIN32_WINDOWS {
						return "Windows 95", nil
					}
				}
			case 10:
				{
					return "Windows 98", nil
				}
			case 90:
				{
					return "Windows Me", nil
				}
			}
		}
	case 5:
		{
			switch o.MinorVersion {
			case 0:
				{
					return "Windows 2000", nil
				}
			case 1:
				{
					return "Windows XP", nil
				}
			case 2:
				// {
				//	r2, err := GetSystemMetrics(SM_SERVERR2)
				//	if err != nil {
				//		return "", err
				//	}
				//
				//	if o.ProductType == byte(VER_NT_WORKSTATION) && u.IsWin64() {
				//		return "Windows XP 64", nil
				//	} else if r2 == 0 {
				//		return "Windows Server 2003", nil
				//	} else if r2 != 0 {
				//		return "Windows Server 2003 R2", nil
				//	}
				// }

			}
		}
	case 6:
		{
			switch o.MinorVersion {
			case 0:
				{
					if o.ProductType == byte(VER_NT_WORKSTATION) {
						return "Windwos Vista", nil
					} else {
						return "Windows Server 2008", nil
					}
				}
			case 1:
				{
					if o.ProductType == byte(VER_NT_WORKSTATION) {
						return "Windwos 7", nil
					} else {
						return "Windows Server 2008 R2", nil
					}
				}
			case 2:
				{
					if o.ProductType == byte(VER_NT_WORKSTATION) {
						return "Windwos 8", nil
					} else {
						return "Windows Server 2012", nil
					}
				}
			}
		}
	}
	return "windows", nil
}
