# WSUSSCN2 CLI - *a CLI for the wsusscn2.cab API*

**Please note that this tool is not affiliated with or created by Microsoft Corporation.**

*Microsoft, Encarta, MSN, and Windows are either registered trademarks or trademarks of Microsoft Corporation in the United States and/or other countries.*

## Goal: Provide a command line interface to Microsoft patch data using the REST API from wsusscn2.cab

## Why

1. For users that need access to Microsoft Windows patch data, but:
  * Don't have access to a WSUS server (or similar enterprise patch tool)
  * Don't want to setup a WSUS server themselves
2. This CLI tool exists for either developers as an example of how to use the wsusscn2.cab API or for users who don't want to program in order to access the patch data

## Requirements (current minimum requirements for golang)

* Windows XP SP3 or newer
* Linux kernel 2.6.23 or newer

## Download

| Windows Release                                                                                                    | MD5                              | SHA1                                     |
|--------------------------------------------------------------------------------------------------------------------|----------------------------------|------------------------------------------|
| [wsusscn2cli v0.3.0](https://github.com/hashauthority/wsusscn2cli/releases/download/v0.3.0/wsusscn2cli-v0.3.0.zip) | ecdb0900f63e623b9bc79d9709055bde | e3693764770636c9ab7de30c1ce126fac4fb8e40 |

* On Windows, run `certutil -hashfile wsusscn2cli-v0.3.0.zip MD5` OR `certutil -hashfile wsusscn2cli-v0.3.0.zip SHA1` to calculate hash of file

| Linux Release                                                                                                         | MD5                              | SHA1                                     |
|-----------------------------------------------------------------------------------------------------------------------|----------------------------------|------------------------------------------|
| [wsusscn2cli v0.3.0](https://github.com/hashauthority/wsusscn2cli/releases/download/v0.3.0/wsusscn2cli-v0.3.0.tar.gz) | 56d2d9b2a3ef49b02156d0d613e84249 | 02f6d485375a0854b61d07f840207a5d4ce1d564 |

* On Linux, run `md5 wsusscn2cli-v0.3.0.tar.gz`, `md5sum wsusscn2cli-v0.3.0.tar.gz`, or `sha1sum wsusscn2cli-v0.3.0.tar.gz` to calculate hash of file

## Getting Started

1. Set API key (Visit https://wsusscn2.cab to create an account and generate an API key).
2. Run `wsusscn2cli setapikey --api_key YOURAPIKEY` to write the API key to wsusscn2cli.json
3. Run `wsusscn2cli listupdates --record_limit 50` and confirm output
4. Run any command with "-q" argument to stop log messages from printing to the screen

## Syntax and examples

Windows patches are "updates" that are released on a typically monthly cadence. Old updates can be superseded by newer updates.

Each row of data represents a unique update for a product. If multiple products have the same update, then each product will have their own row of data.

```
> wsusscn2cli -h
NAME:
   wsusscn2cli - wsusscn2.cab integration

USAGE:
   wsusscn2cli [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
     listclassification  List all classifications
     listproduct         List all products
     listproductfamily   List all product families
     listupdate          List updates
     listsupersede       List supersession updates
     setapikey           Set API key for repeated usage
     help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   (c) 2018 Hash Authority, LLC
```

### **```wsusscn2cli listupdate```**

```
> wsusscn2cli listupdate -h
NAME:
   wsusscn2cli listupdate - List updates

USAGE:
   wsusscn2cli listupdate [command options] [arguments...]

OPTIONS:
   --api_key value, -a value            API key (required if not using config file)
   --debug, -d                          Output debug level logging
   --insecure, -k                       Do not verify server's SSL cert
   --quiet, -q                          Do not log to screen
   --count_only                         Only print number of records
   --product_title value                Name of product.
   --update_uid value                   Update Uid.
   --update_title value                 Update Title.
   --kb value                           Update KB.
   --update_type value                  Update Type.
   --product_family_title value         Product Family Title.
   --classification_title value         Classification Title.
   --msrc_severity value                MSRC Severity.
   --arch value                         Architecture.
   --is_superseded value                Is Superseded.
   --is_bundled value                   Is Bundled.
   --is_public value                    Is Public.
   --is_beta value                      Is Beta.
   --update_creation_date_after value   Updates created after this date [YYYY-MM-DD] (exclusive).
   --update_creation_date_before value  Updates created before this date [YYYY-MM-DD] (exclusive).
   --update_creation_date_on value      Updates created on this date [YYYY-MM-DD].
   --columns value                      Restrict output to listed columns.
   --limit value                        Number of records per page. (default: 1000)
   --offset value                       Number of records to skip. (default: 0)
   --record_limit value                 Max number of records to return. (default: 20000)
```

Definition: List updates. Multiple values can be passed by repeating the argument. Multiple arguments for the same filter are ORed together. Filters are ANDed together. Use --columns to reduce the fields in output.

* "update_uid": Unique identifier for an update. Combine this with product_title to get a single row of data.
* "kb": KB number
* "update_title": Update Title
* "update_type": Update Type
* "product_family_title": Product Family Title
* "classification_title": Classification Title
* "count_only": List the number of records returned
* "product_title": OS or Application name
* "msrc_severity": Severity rating of patch by Microsoft
* "is_superseded": Indicates if this is superseded by another update. Values allowed can be 0/1, t/f, true/false, True/False
* "is_public": Indicates if this update is released to the public. Values allowed can be 0/1, t/f, true/false, True/False
* "is_beta": Indicates if this update is in beta. Values allowed can be 0/1, t/f, true/false, True/False
* "is_bundled": Indicates if this update is included in another update.

Example of updates for Windows 7:
```
> wsusscn2cli listupdate --product_title "Windows 7" --record_limit 5 --columns "update_uid, kb, update_title, product_title"
"UpdateUid","Kb","UpdateTitle","ProductTitle","ProductFamilyTitle"
"E302AE72-7CF8-48CE-9B19-A9E28E197280","4092946","Cumulative Security Update for Internet Explorer 11 for Windows 7 (KB4092946)","Windows 7","Windows"
"71E638FE-A799-4166-9C75-56A8D1263C2E","4092946","Cumulative Security Update for Internet Explorer 11 for Windows 7 for x64-based Systems (KB4092946)","Windows 7","Windows"
"D1FDFCFA-0E2E-4EF9-AAD2-F97E1EA108D1","890830","Windows Malicious Software Removal Tool x64 - April 2018 (KB890830)","Windows 7","Windows"
"E364176B-CF12-4880-B745-D25BC2603027","890830","Windows Malicious Software Removal Tool - April 2018 (KB890830)","Windows 7","Windows"
"520C9A8F-BC91-42C2-9C5F-5424F80E8349","4093118","2018-04 Security Monthly Quality Rollup for Windows 7 for x86-based Systems (KB4093118)","Windows 7","Windows"
```

Example of updates for Windows 7 or Windows 10:
```
> wsusscn2cli listupdate --product_title "Windows 7" --product_title "Windows 10" --record_limit 5 --columns "update_uid, kb, update_title, product_title"
"UpdateUid","Kb","UpdateTitle","ProductTitle"
"E302AE72-7CF8-48CE-9B19-A9E28E197280","4092946","Cumulative Security Update for Internet Explorer 11 for Windows 7 (KB4092946)","Windows 7"
"71E638FE-A799-4166-9C75-56A8D1263C2E","4092946","Cumulative Security Update for Internet Explorer 11 for Windows 7 for x64-based Systems (KB4092946)","Windows 7"
"D1FDFCFA-0E2E-4EF9-AAD2-F97E1EA108D1","890830","Windows Malicious Software Removal Tool x64 - April 2018 (KB890830)","Windows 7"
"E364176B-CF12-4880-B745-D25BC2603027","890830","Windows Malicious Software Removal Tool - April 2018 (KB890830)","Windows 7"
"A32121DF-DD27-4B32-9ECC-927E9915E083","890830","Windows Malicious Software Removal Tool - April 2018 (KB890830)","Windows 10"
```

Example of critical severity updates for Windows 7 and Windows 10:
```
> wsusscn2cli listupdate --product_title "Windows 7" --product_title "Windows 10" --msrc_severity "Critical" --record_limit 5 --columns "update_uid, kb, update_title, product_title"
"UpdateUid","Kb","UpdateTitle","ProductTitle"
"E302AE72-7CF8-48CE-9B19-A9E28E197280","4092946","Cumulative Security Update for Internet Explorer 11 for Windows 7 (KB4092946)","Windows 7"
"71E638FE-A799-4166-9C75-56A8D1263C2E","4092946","Cumulative Security Update for Internet Explorer 11 for Windows 7 for x64-based Systems (KB4092946)","Windows 7"
"520C9A8F-BC91-42C2-9C5F-5424F80E8349","4093118","2018-04 Security Monthly Quality Rollup for Windows 7 for x86-based Systems (KB4093118)","Windows 7"
"954A4DC7-6623-4156-95D1-AE1296052BF6","4093108","2018-04 Security Only Quality Update for Windows 7 for x86-based Systems (KB4093108)","Windows 7"
"C25445CD-6E70-42FB-965F-7650E629CC42","4093108","2018-04 Security Only Quality Update for Windows 7 for x64-based Systems (KB4093108)","Windows 7"
```

Example of kb search:
```
> wsusscn2cli listupdate --kb 2923392 --columns "update_title, kb, update_creation_date, product_title"
"UpdateTitle","Kb","UpdateCreationDate","ProductTitle"
"Security Update for Windows Vista for x64-based Systems (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Vista"
"Security Update for Windows Vista (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Vista"
"Security Update for Windows Server 2003 (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Server 2003"
"Security Update for Windows Server 2003 (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Server 2003, Datacenter Edition"
"Security Update for Windows Server 2003 x64 Edition (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Server 2003"
"Security Update for Windows Server 2003 x64 Edition (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Server 2003, Datacenter Edition"
"Security Update for Windows Server 2003 for Itanium-based Systems (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Server 2003"
"Security Update for Windows Server 2003 for Itanium-based Systems (KB2923392)","2923392","2014-03-11T17:00:00Z","Windows Server 2003, Datacenter Edition"
```

Example of counting important updates for "Windows 7"
```
> wsusscn2cli listupdate --product_title "Windows 7" --count_only --msrc_severity "Important"
Number of records: 466
```

### **```wsusscn2cli listsupersede```**

```
> wsusscn2cli.exe listsupersede -h
NAME:
   wsusscn2cli.exe listsupersede - List supersession updates

USAGE:
   wsusscn2cli.exe listsupersede [command options] [arguments...]

OPTIONS:
   --api_key value, -a value            API key (required if not using config file)
   --debug, -d                          Output debug level logging
   --quiet, -q                          Do not log to screen
   --product_title value                Name of product.
   --update_uid value                   Update Uid.
   --update_title value                 Update Title.
   --kb value                           Update KB.
   --update_type value                  Update Type.
   --product_family_title value         Product Family Title.
   --classification_title value         Classification Title.
   --msrc_severity value                MSRC Severity.
   --is_superseded value                Is Superseded.
   --is_bundled value                   Is Bundled.
   --is_public value                    Is Public.
   --is_beta value                      Is Beta.
   --update_creation_date_after value   Updates created after this date [YYYY-MM-DD] (exclusive).
   --update_creation_date_before value  Updates created before this date [YYYY-MM-DD] (exclusive).
   --update_creation_date_on value      Updates created on this date [YYYY-MM-DD].
   --limit value                        Number of records per page. (default: 1000)
   --offset value                       Number of records to skip. (default: 0)
   --record_limit value                 Max number of records to return. (default: 20000)
```

Definition: List superseded updates. The listed update is the *latest* update in the supersession chain. That is, if A is superseded by B and B is superseded by C, then this command will show that A is superseded by C.

* "update_uid": Unique identifier for an update. Combine this with product_title to get a single row of data.
* "kb": KB number
* "update_title": Update Title
* "update_type": Update Type
* "product_family_title": Product Family Title
* "classification_title": Classification Title
* "product_title": OS or Application name
* "msrc_severity": Severity rating of patch by Microsoft
* "is_superseded": Indicates if this is superseded by another update. Values allowed can be 0/1, t/f, true/false, True/False
* "is_public": Indicates if this update is released to the public. Values allowed can be 0/1, t/f, true/false, True/False
* "is_beta": Indicates if this update is in beta. Values allowed can be 0/1, t/f, true/false, True/False
* "is_bundled": Indicates if this update is included in another update.

Example of superseded by kb:
```
> wsusscn2cli listsupersede --kb 4025339
"UpdateUid","UpdateTitle","UpdateCreationDate","ProductTitle","IsSuperseded","SuperUpdateUid","SuperTitle","SuperCreationDate","SuperProductTitle","SuperIsSuperseded"
"84C1A786-79D4-4D94-9A75-FC900083816F","2017-07 Cumulative Update for Windows 10 Version 1607 for x86-based Systems (KB4025339)","2017-07-11T17:00:03Z","Windows 10","true","073AA939-731E-464C-B64E-F6241C4D9A86","2018-06 Cumulative Update for Windows 10 Version 1607 for x86-based Systems (KB4284880)","2018-06-12T17:00:05Z","Windows 10","false"
"84C1A786-79D4-4D94-9A75-FC900083816F","2017-07 Cumulative Update for Windows 10 Version 1607 for x86-based Systems (KB4025339)","2017-07-11T17:00:03Z","Windows 10","true","073AA939-731E-464C-B64E-F6241C4D9A86","2018-06 Cumulative Update for Windows 10 Version 1607 for x86-based Systems (KB4284880)","2018-06-12T17:00:05Z","Windows 10 LTSB","false"
"A44D500F-FC5D-4FE4-90CB-991568E9CB58","2017-07 Cumulative Update for Windows 10 Version 1607 for x64-based Systems (KB4025339)","2017-07-11T17:00:03Z","Windows 10","true","3105F320-7D78-4034-A86C-03B4F9352480","2018-06 Cumulative Update for Windows 10 Version 1607 for x64-based Systems (KB4284880)","2018-06-12T17:00:05Z","Windows 10","false"
"A44D500F-FC5D-4FE4-90CB-991568E9CB58","2017-07 Cumulative Update for Windows 10 Version 1607 for x64-based Systems (KB4025339)","2017-07-11T17:00:03Z","Windows 10","true","3105F320-7D78-4034-A86C-03B4F9352480","2018-06 Cumulative Update for Windows 10 Version 1607 for x64-based Systems (KB4284880)","2018-06-12T17:00:05Z","Windows 10 LTSB","false"
"D6677D54-CE7A-4774-A696-84DE34EFF033","2017-07 Cumulative Update for Windows Server 2016 for x64-based Systems (KB4025339)","2017-07-11T17:00:03Z","Windows Server 2016","true","FA8B8608-4925-4C9B-871F-A3E5D0B082FA","2018-06 Cumulative Update for Windows Server 2016 for x64-based Systems (KB4284880)","2018-06-12T17:00:05Z","Windows Server 2016","false"
```

### **```wsusscn2cli listcve```**

```
> wsusscn2cli.exe listcve -h
NAME:
   wsusscn2cli.exe listcve - List all CVEs

USAGE:
   wsusscn2cli.exe listcve [command options] [arguments...]

OPTIONS:
   --api_key value, -a value      API key (required if not using config file)
   --debug, -d                    Output debug level logging
   -k, --insecure                 Do not verify server's SSL cert
   --quiet, -q                    Do not log to screen
   --cve value                    CVE number (Ex., CVE-2018-0001).
   --cvssv3_base_score value      CVSS v3 Base Score (Range 1-10). Range allowed (Ex., 7.1-10.0)
   --cvssv3_temporal_score value  CVSS v3 Temporal Score (Range 1-10). Range allowed (Ex., 7.1-10.0)
   --product_title value          Name of product.
   --update_uid value             Update Uid.
   --update_title value           Update Title.
   --kb value                     Update KB.
   --product_family_title value   Product Family Title.
   --classification_title value   Classification Title.
   --msrc_severity value          MSRC Severity.
   --arch value                   Architecture.
   --is_superseded value          Is Superseded.
   --is_in_file value             Is in file (is in the current wsusscn2.cab file).
   --limit value                  Number of records per page. (default: 1000)
   --offset value                 Number of records to skip. (default: 0)
   --record_limit value           Max number of records to return. (default: 20000)
```

### **```wsusscn2cli listclassification```**

```
> wsusscn2cli listclassification -h
NAME:
   wsusscn2cli listclassification - List all classifications

USAGE:
   wsusscn2cli listclassification [command options] [arguments...]

OPTIONS:
   --debug, -d                Output debug level logging
   --api_key value, -a value  API key (required if not using config file)
```

Definition: Display list of classifications

Example:
```
> wsusscn2cli listclassification
"ClassificationUid","ClassificationRevision","ClassificationTitle"
"68C5B0A3-D1A6-4553-AE49-01D3A7827828","9","Service Packs"
"28BC880E-0592-4CBF-8F95-C79B17911D5F","8","Update Rollups"
"CD5FFD1E-E932-4E3A-BF74-18BF0B1BBD83","7","Updates"
"E6CF1350-C01B-414D-A61F-263D14D133B4","6","Critical Updates"
"0FA1201D-4330-4FA8-8AE9-B877473B6441","5","Security Updates"
```

### **```wsusscn2cli listproduct```**

```
> wsusscn2cli listproduct -h
NAME:
   wsusscn2cli listproduct - List all products

USAGE:
   wsusscn2cli listproduct [command options] [arguments...]

OPTIONS:
   --debug, -d                Output debug level logging
   --api_key value, -a value  API key (required if not using config file)
```

Definition: Display list of products

Example:
```
> wsusscn2cli listproduct
"ProductUid","ProductRevision","ProductTitle"
"7FF1D901-FD38-441B-AABA-36D7B0EBF264","25766777","Azure File Sync agent updates for Windows Server 2016"
"FB08C71C-DBE9-40AB-8302-FB0231B1C814","25766776","Azure File Sync agent updates for Windows Server 2012 R2"
"A3C2375D-0C8A-42F9-BCE0-28333E198407","25629036","Windows 10"
"CA6616AA-6310-4C2D-A6BF-CAE700B85E86","25436193","Microsoft SQL Server 2017"
"589DB546-7849-47F5-BBC0-1F66CF12F5C2","24677545","Windows 8 Embedded"
[snip]
```

### **```wsusscn2cli listproductfamily```**

```
> wsusscn2cli listproductfamily -h
NAME:
   wsusscn2cli listproductfamily - List all product families

USAGE:
   wsusscn2cli listproductfamily [command options] [arguments...]

OPTIONS:
   --debug, -d                Output debug level logging
   --api_key value, -a value  API key (required if not using config file)
   --quiet, -q                Do not log to screen
```

Definition: Display list of product families

Example:
```
> wsusscn2cli listproductfamily
"ProductFamilyUid","ProductFamilyRevision","ProductFamilyTitle"
"9D6F2556-534F-047E-5EC9-91BF0DA81A75","25692907","Azure File Sync"
"0DBC842C-730F-4361-8811-1B048F11C09B","21374218","Microsoft Dynamics CRM"
"2E97A7D7-8256-58EF-2FB6-48CBACDB603D","21230345","Microsoft Advanced Threat Analytics"
"8FDC8B60-9E7C-4275-8668-198F89A64DF6","17931119","Skype for Business"
"64DA29AF-92B5-C36A-FAB2-682350A63C2F","16688976","Microsoft Monitoring Agent (MMA)"
"4756F399-B049-8E6E-94E9-FF63D0E236A7","13832531","ASP.NET Web and Data Frameworks"
[snip]
```

### **```wsusscn2cli setapikey```**

```
> wsusscn2cli setapikey -h
NAME:
   wsusscn2cli setapikey - Set API key for repeated usage

USAGE:
   wsusscn2cli setapikey [command options] [arguments...]

OPTIONS:
   --debug, -d                Output debug level logging
   --api_key value, -a value  Authentication to API
```

Definition: Set the API key used for authentication to wsusscn2.cab API

Example:
```
> wsusscn2cli setapikey --api_key e685304f4c1d57d7bd7a59ab9c159e9d
```

## Version history
* **0.1.0** (2018-04-09) - Internal release only.
* **0.1.1** (2018-04-10) - Internal release only.
* **0.1.2** (2018-04-12) - Internal release only.
* **0.1.3** (2018-04-25) - Internal release only.
* **0.1.4** (2018-06-11) - Initial public release. Release of wsusscn2cli binary with listupdate, listclassification, listproduct, and listproductfamily commands
* **0.1.5** (unreleased) - Added listsupersede command, fixed bug with update_creation_date_on argument, and added quiet argument to stop logging to the screen
* **0.2.0** (2018-09-30) - Updated endpoint to api.wsusscn2.cab. Note that all previous versions will no longer work since the root domain is now a web page.
* **0.3.0** (2018-10-12) - Added listcve command. Added --insecure switch to ignore server ssl cert verification (should not be required for most environments).

## License

wsusscn2cli is licensed under the [MIT](http://www.opensource.org/licenses/mit-license.php). However, in order for this tool to work, you must have an existing license to the wsusscn2.cab API and a valid API key from [https://wsusscn2.cab](https://wsusscn2.cab)

## Other libraries used by wsusscn2cli

* [urfave/cli](https://github.com/urfave/cli) *(MIT License)*
