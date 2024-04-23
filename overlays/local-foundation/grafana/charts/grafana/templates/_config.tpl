{{- define "grafana.secretsData" -}}
{{- if and (not .Values.env.GF_SECURITY_DISABLE_INITIAL_ADMIN_CREATION) (not .Values.admin.existingSecret) (not .Values.env.GF_SECURITY_ADMIN_PASSWORD__FILE) (not .Values.env.GF_SECURITY_ADMIN_PASSWORD) }}
admin-user: {{ .Values.adminUser | b64enc | quote }}
{{- if .Values.adminPassword }}
admin-password: {{ .Values.adminPassword | b64enc | quote }}
{{- else }}
admin-password: {{ include "grafana.password" . }}
{{- end }}
{{- end }}
{{- if not .Values.ldap.existingSecret }}
ldap-toml: {{ tpl .Values.ldap.config $ | b64enc | quote }}
{{- end }}
{{- end -}}