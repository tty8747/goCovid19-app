{{- define "gocovid.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "namespace" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "dbString" -}}
{{- printf "%s:%s@tcp(%s:%d)/%s" .Values.dbSettings.user .Values.dbSettings.password .Values.dbSettings.endpoint (.Values.dbSettings.port | int) .Values.dbSettings.name | b64enc }}
{{- end }}

{{- define "apiConfig" -}}
{{ tpl (.Files.Get "conf/config.yml") . | b64enc }}
{{- end }}

{{- define "gooseImage" -}}
{{- printf "%s:%s" .Values.goose.image.repository .Values.goose.image.tag }}
{{- end }}

{{- define "apiImage" -}}
{{- printf "%s:%s" .Values.api.image.repository .Values.api.image.tag }}
{{- end }}

{{- define "frontImage" -}}
{{- printf "%s:%s" .Values.front.image.repository .Values.front.image.tag }}
{{- end }}
