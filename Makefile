all: filetpl paramtpl parseDate repeat md5 toPrettyJson
.PHONY: filetpl
filetpl:
	@echo '{"FirstName": "John", "LastName": "Doe"}' | go run . example.tmpl

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
