{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "mongodb.fullname" -}}
{{- printf "%s-%s" .Release.Name "mongodb" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Render image reference
*/}}
{{- define "monocular.image" -}}
{{ .registry }}/{{ .repository }}:{{ .tag }}
{{- end -}}

{{/*
Sync job pod template
*/}}
{{- define "monocular.sync.podTemplate" -}}
{{- $repo := index . 0 -}}
{{- $global := index . 1 -}}
metadata:
  labels:
    monocular.helm.sh/repo-name: {{ $repo.name }}
    app: {{ template "fullname" $global }}
    release: "{{ $global.Release.Name }}"
spec:
  restartPolicy: OnFailure
  containers:
  - name: sync
    image: {{ template "monocular.image" $global.Values.sync.image }}
    args:
    - sync
    - --user-agent-comment=monocular/{{ $global.Chart.AppVersion }}
    {{- if and $global.Values.global.mongoUrl (not $global.Values.mongodb.enabled) }}
    - --mongo-url={{ $global.Values.global.mongoUrl }}
    {{- else }}
    - --mongo-url={{ template "mongodb.fullname" $global }}
    - --mongo-user=root
    {{- end }}
    - {{ $repo.name }}
    - {{ $repo.url }}
    command:
    - /chart-repo
    {{- if $global.Values.mongodb.enabled }}
    env:
    - name: HTTP_PROXY
      value: {{ $global.Values.sync.httpProxy }}
    - name: HTTPS_PROXY
      value: {{ $global.Values.sync.httpsProxy }}
    - name: NO_PROXY
      value: {{ $global.Values.sync.noProxy }}
    - name: MONGO_PASSWORD
      valueFrom:
        secretKeyRef:
          key: mongodb-root-password
          name: {{ template "mongodb.fullname" $global }}
    {{- end }}
    resources:
{{ toYaml $global.Values.sync.resources | indent 6 }}
{{- with $global.Values.sync.nodeSelector }}
  nodeSelector:
{{ toYaml . | indent 4 }}
{{- end }}
{{- with $global.Values.sync.affinity }}
  affinity:
{{ toYaml . | indent 4 }}
{{- end }}
{{- with $global.Values.sync.tolerations }}
  tolerations:
{{ toYaml . | indent 4 }}
{{- end }}
{{- end -}}
