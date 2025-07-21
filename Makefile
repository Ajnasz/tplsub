build:
	@go build -o tplsub
test: nodata filetpl datafile paramtpl parseDate repeat md5 toPrettyJson stringHelpers mathHelpers floatMathHelpers dateHelpers collectionHelpers conditionalHelpers fileHelpers envHelpers hashingHelpers

.PHONY: nodata
nodata:
	@go run . -t "hello world"
	@echo

.PHONY: filetpl
filetpl:
	@echo '{"FirstName": "John", "LastName": "Doe"}' | go run . examples/example.tmpl

.PHONY: datafile
datafile:
	@go run . -t 'Data file: {{ .FirstName }} {{ .LastName }}' examples/data.json
	@echo

.PHONY: paramtpl
paramtpl:
	@echo '{"FirstName": "John", "LastName": "Doe"}' | go run . --template 'Param tpl: {{ .FirstName }} {{ .LastName | lower }}'
	@echo

.PHONY: parseDate
parseDate:
	@echo '{"time": "2025-07-20 17:17:00 CEST"}' | go run . --template 'min: {{ .time | parseDate "2006-01-02 15:04:05 MST" | year }}'
	@echo

.PHONY: repeat
repeat:
	@echo '{"FirstName": "John", "LastName": "Doe"}' | go run . --template 'Repeat: {{ repeat .FirstName 3 }} {{ repeat .LastName  2 }}'
	@echo

.PHONY: md5
md5:
	@echo '{"FirstName": "John", "LastName": "Doe"}' | go run . --template 'MD5: {{ md5 .FirstName }} {{ md5 .LastName }}'
	@echo

.PHONY: toPrettyJSON
toPrettyJson:
	@echo '{"FirstName": "John", "LastName": "Doe"}' | go run . --template '{{ . | toPrettyJSON }}'
	@echo

.PHONY: stringHelpers
stringHelpers:
	@echo 'String manipulation examples:'
	@echo '{"text": "Hello World"}' | go run . --template 'Upper: {{ .text | upper }}'
	@echo
	@echo '{"text": "Hello World"}' | go run . --template 'Lower: {{ .text | lower }}'
	@echo
	@echo '{"text": "  hello  "}' | go run . --template 'Trim: "{{ .text | trim }}"'
	@echo
	@echo '{"text": "Hello World"}' | go run . --template 'Replace: "{{ .text | replace "World" "Gopher" }}"'
	@echo
	@echo '{"text": "hello,world,test"}' | go run . --template 'Split: {{ split "," .text }}'
	@echo
	@echo '{"items": ["a", "b", "c"]}' | go run . --template 'Join: {{ join "-" .items }}'
	@echo
	@echo '{"text": "hello world"}' | go run . --template 'Contains "world": {{ contains "world" .text }}'
	@echo
	@echo '{"text": "hello world"}' | go run . --template 'Has prefix "hello": {{ .text | hasPrefix "hello" }}'
	@echo
	@echo '{"text": "hello world"}' | go run . --template 'Has suffix "world": {{ .text | hasSuffix "world" }}'
	@echo
	@echo

.PHONY: mathHelpers
mathHelpers:
	@echo 'Math operation examples:'
	@echo
	@echo '{"a": 10, "b": 3}' | go run . --template 'Add: {{ add .a .b }}'
	@echo
	@echo '{"a": 10, "b": 3}' | go run . --template 'Sub: {{ sub .a .b }}'
	@echo
	@echo '{"a": 10, "b": 3}' | go run . --template 'Mul: {{ mul .a .b }}'
	@echo
	@echo '{"a": 10, "b": 3}' | go run . --template 'Div: {{ div .a .b }}'
	@echo
	@echo '{"a": 10, "b": 3}' | go run . --template 'Mod: {{ mod .a .b }}'
	@echo
	@echo

.PHONY: floatMathHelpers
floatMathHelpers:
	@echo 'Float math operation examples:'
	@echo
	@echo '{"a": 10.5, "b": 3.2}' | go run . --template 'Add float: {{ addf .a .b }}'
	@echo
	@echo '{"a": 10.5, "b": 3.2}' | go run . --template 'Sub float: {{ subf .a .b }}'
	@echo
	@echo '{"a": 10.5, "b": 3.2}' | go run . --template 'Mul float: {{ mulf .a .b }}'
	@echo
	@echo '{"a": 10.5, "b": 3.2}' | go run . --template 'Div float: {{ divf .a .b }}'
	@echo
	@echo '{"a": "15.75", "b": 4}' | go run . --template 'Mixed types: {{ addf .a .b }} (string + int)'
	@echo
	@echo '{"value": 42}' | go run . --template 'To float: {{ toFloat .value }} (converted to float)'
	@echo
	@echo '{"a": 22, "b": 7}' | go run . --template 'Precise division: {{ divf .a .b }} vs {{ div .a .b }}'
	@echo
	@echo

.PHONY: dateHelpers
dateHelpers:
	@echo 'Date/time helper examples:'
	@echo
	@echo '{}' | go run . --template 'Now: {{ now }}'
	@echo
	@echo '{}' | go run . --template 'Formatted: {{ now | formatDate "2006-01-02 15:04:05" }}'
	@echo
	@echo '{}' | go run . --template 'Year: {{ now | year }}'
	@echo
	@echo '{"date": "2023-12-25"}' | go run . --template 'Parsed: {{ parseDate "2006-01-02" .date }}'
	@echo
	@echo

.PHONY: collectionHelpers
collectionHelpers:
	@echo 'Collection helper examples:'
	@echo
	@echo '{"items": ["first", "second", "third"]}' | go run . --template 'Length: {{ len .items }}'
	@echo
	@echo '{"items": ["first", "second", "third"]}' | go run . --template 'First: {{ first .items }}'
	@echo
	@echo '{"items": ["first", "second", "third"]}' | go run . --template 'Last: {{ last .items }}'
	@echo
	@echo '{}' | go run . --template 'Sequence: {{ range seq 1 5 }}{{ . }} {{ end }}'
	@echo
	@echo

.PHONY: conditionalHelpers
conditionalHelpers:
	@echo 'Conditional helper examples:'
	@echo
	@echo '{"name": ""}' | go run . --template 'Default: {{ default "Anonymous" .name }}'
	@echo
	@echo '{"name": "John"}' | go run . --template 'Default: {{ default "Anonymous" .name }}'
	@echo
	@echo '{"name": ""}' | go run . --template 'Empty check: {{ if empty .name }}Name is empty{{ else }}Name: {{ .name }}{{ end }}'
	@echo
	@echo

.PHONY: fileHelpers
fileHelpers:
	@echo 'File path helper examples:'
	@echo
	@echo '{"path": "/home/user/document.txt"}' | go run . --template 'Basename: {{ basename .path }}'
	@echo
	@echo '{"path": "/home/user/document.txt"}' | go run . --template 'Dirname: {{ dirname .path }}'
	@echo
	@echo '{"path": "/home/user/document.txt"}' | go run . --template 'Extension: {{ ext .path }}'
	@echo
	@echo '{"dir": "/home", "subdir": "user", "file": "doc.txt"}' | go run . --template 'Join: {{ pathjoin .dir .subdir .file }}'
	@echo
	@echo

.PHONY: envHelpers
envHelpers:
	@echo 'Environment helper examples:'
	@echo
	@echo '{}' | go run . --template 'Home: {{ env "HOME" }}'
	@echo
	@echo '{}' | go run . --template 'User: {{ env "USER" }}'
	@echo
	@echo

.PHONY: hashingHelpers
hashingHelpers:
	@echo 'Hashing and encoding examples:'
	@echo
	@echo '{"text": "hello world"}' | go run . --template 'MD5: {{ md5 .text }}'
	@echo
	@echo '{"text": "hello world"}' | go run . --template 'SHA1: {{ sha1 .text }}'
	@echo
	@echo '{"text": "hello world"}' | go run . --template 'SHA256: {{ sha256 .text }}'
	@echo
	@echo '{"text": "hello world"}' | go run . --template 'Base64: {{ base64Encode .text }}'
	@echo
	@echo '{"encoded": "aGVsbG8gd29ybGQ="}' | go run . --template 'Decoded: {{ base64Decode .encoded }}'
	@echo
	@echo
