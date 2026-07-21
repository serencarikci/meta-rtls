import { FormEvent, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  Alert,
  Box,
  Button,
  Chip,
  MenuItem,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from '@mui/material'
import { api } from '../api/client'
import { useAuthStore } from '../store/auth'

type Definition = {
  id: string
  entityType: string
  code: string
  name: string
  description: string
  currentVersion: number
  status: string
}

type SchemaVersion = {
  id: string
  versionNo: number
  changelog: string
  isCurrent: boolean
}

type Field = {
  id: string
  fieldKey: string
  label: string
  dataType: string
  isRequired: boolean
  minValue?: number
  maxValue?: number
  enumValues?: string[]
}

type TenantFeature = {
  code: string
  name: string
  enabled: boolean
  category: string
}

type ValidateResponse = {
  valid: boolean
  errors: string[]
}

const ENTITY_TYPES = ['ASSET', 'PERSON', 'TAG', 'ZONE', 'DEVICE']
const DATA_TYPES = ['STRING', 'NUMBER', 'BOOLEAN', 'ENUM', 'DATE', 'JSON']

export default function MetadataPage() {
  const token = useAuthStore((s) => s.token)
  const queryClient = useQueryClient()
  const [selectedId, setSelectedId] = useState<string>('')

  const [defForm, setDefForm] = useState({
    entityType: 'ASSET',
    code: '',
    name: '',
    description: '',
  })
  const [fieldForm, setFieldForm] = useState({
    fieldKey: '',
    label: '',
    dataType: 'STRING',
    isRequired: true,
    enumValuesText: '',
  })
  const [validateJson, setValidateJson] = useState(
    '{\n  "loadCapacity": 1500,\n  "batteryLevel": 82,\n  "department": "Warehouse"\n}',
  )
  const [validateResult, setValidateResult] = useState<ValidateResponse | null>(null)

  const definitions = useQuery({
    queryKey: ['metadata-definitions'],
    queryFn: () => api<Definition[]>('/api/v1/metadata/definitions', {}, token),
  })

  const features = useQuery({
    queryKey: ['metadata-features'],
    queryFn: () => api<TenantFeature[]>('/api/v1/metadata/features', {}, token),
  })

  const activeId = selectedId || definitions.data?.[0]?.id || ''

  const versions = useQuery({
    queryKey: ['metadata-versions', activeId],
    enabled: !!activeId,
    queryFn: () =>
      api<SchemaVersion[]>(`/api/v1/metadata/definitions/${activeId}/versions`, {}, token),
  })

  const currentVersion = useMemo(
    () => versions.data?.find((v) => v.isCurrent) ?? versions.data?.[versions.data.length - 1],
    [versions.data],
  )

  const fields = useQuery({
    queryKey: ['metadata-fields', currentVersion?.id],
    enabled: !!currentVersion?.id,
    queryFn: () =>
      api<Field[]>(`/api/v1/metadata/versions/${currentVersion!.id}/fields`, {}, token),
  })

  const createDefinition = useMutation({
    mutationFn: () =>
      api<Definition>(
        '/api/v1/metadata/definitions',
        { method: 'POST', body: JSON.stringify(defForm) },
        token,
      ),
    onSuccess: (item) => {
      queryClient.invalidateQueries({ queryKey: ['metadata-definitions'] })
      setSelectedId(item.id)
      setDefForm({ entityType: 'ASSET', code: '', name: '', description: '' })
    },
  })

  const createField = useMutation({
    mutationFn: () => {
      const body = {
        fieldKey: fieldForm.fieldKey,
        label: fieldForm.label,
        dataType: fieldForm.dataType,
        isRequired: fieldForm.isRequired,
        enumValues:
          fieldForm.dataType === 'ENUM'
            ? fieldForm.enumValuesText
                .split(',')
                .map((x) => x.trim())
                .filter(Boolean)
            : [],
      }
      return api<Field>(
        `/api/v1/metadata/versions/${currentVersion!.id}/fields`,
        { method: 'POST', body: JSON.stringify(body) },
        token,
      )
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metadata-fields', currentVersion?.id] })
      setFieldForm({
        fieldKey: '',
        label: '',
        dataType: 'STRING',
        isRequired: true,
        enumValuesText: '',
      })
    },
  })

  const createVersion = useMutation({
    mutationFn: () =>
      api<SchemaVersion>(
        `/api/v1/metadata/definitions/${activeId}/versions`,
        { method: 'POST', body: JSON.stringify({ changelog: 'manual bump' }) },
        token,
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metadata-versions', activeId] })
      queryClient.invalidateQueries({ queryKey: ['metadata-definitions'] })
    },
  })

  async function onCreateDefinition(e: FormEvent) {
    e.preventDefault()
    createDefinition.mutate()
  }

  async function onCreateField(e: FormEvent) {
    e.preventDefault()
    createField.mutate()
  }

  async function onValidate() {
    try {
      const values = JSON.parse(validateJson)
      const result = await api<ValidateResponse>(
        '/api/v1/metadata/validate',
        {
          method: 'POST',
          body: JSON.stringify({ definitionId: activeId, values }),
        },
        token,
      )
      setValidateResult(result)
    } catch (err) {
      setValidateResult({
        valid: false,
        errors: [err instanceof Error ? err.message : 'validate failed'],
      })
    }
  }

  return (
    <Stack spacing={3}>
      <div>
        <Typography variant="h4">Metadata engine</Typography>
        <Typography color="text.secondary">
          Define fields per tenant. Validate values without code changes.
        </Typography>
      </div>

      <Paper sx={{ p: 2, border: '1px solid #D5DEE8' }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Tenant features
        </Typography>
        <Stack direction="row" gap={1} flexWrap="wrap">
          {(features.data ?? []).map((feature) => (
            <Chip
              key={feature.code}
              label={`${feature.code}${feature.enabled ? '' : ' (off)'}`}
              color={feature.enabled ? 'success' : 'default'}
              variant="outlined"
              size="small"
            />
          ))}
        </Stack>
      </Paper>

      <Paper sx={{ p: 2, border: '1px solid #D5DEE8' }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Definitions
        </Typography>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Code</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Entity</TableCell>
              <TableCell>Version</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>
            {(definitions.data ?? []).map((item) => (
              <TableRow key={item.id} selected={item.id === activeId}>
                <TableCell>{item.code}</TableCell>
                <TableCell>{item.name}</TableCell>
                <TableCell>{item.entityType}</TableCell>
                <TableCell>v{item.currentVersion}</TableCell>
                <TableCell>
                  <Button size="small" onClick={() => setSelectedId(item.id)}>
                    Open
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <Box component="form" onSubmit={onCreateDefinition} sx={{ mt: 2 }}>
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={1}>
            <TextField
              select
              label="Entity"
              size="small"
              value={defForm.entityType}
              onChange={(e) => setDefForm({ ...defForm, entityType: e.target.value })}
              sx={{ minWidth: 140 }}
            >
              {ENTITY_TYPES.map((item) => (
                <MenuItem key={item} value={item}>
                  {item}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              label="Code"
              size="small"
              value={defForm.code}
              onChange={(e) => setDefForm({ ...defForm, code: e.target.value })}
              required
            />
            <TextField
              label="Name"
              size="small"
              value={defForm.name}
              onChange={(e) => setDefForm({ ...defForm, name: e.target.value })}
              required
            />
            <Button type="submit" variant="contained" disabled={createDefinition.isPending}>
              Add definition
            </Button>
          </Stack>
        </Box>
      </Paper>

      {activeId && (
        <Paper sx={{ p: 2, border: '1px solid #D5DEE8' }}>
          <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
            <Typography variant="h6">
              Current fields {currentVersion ? `(v${currentVersion.versionNo})` : ''}
            </Typography>
            <Button size="small" onClick={() => createVersion.mutate()} disabled={!activeId}>
              New schema version
            </Button>
          </Stack>

          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Key</TableCell>
                <TableCell>Label</TableCell>
                <TableCell>Type</TableCell>
                <TableCell>Required</TableCell>
                <TableCell>Rules</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {(fields.data ?? []).map((field) => (
                <TableRow key={field.id}>
                  <TableCell>{field.fieldKey}</TableCell>
                  <TableCell>{field.label}</TableCell>
                  <TableCell>{field.dataType}</TableCell>
                  <TableCell>{field.isRequired ? 'yes' : 'no'}</TableCell>
                  <TableCell>
                    {field.minValue != null || field.maxValue != null
                      ? `${field.minValue ?? '-'}..${field.maxValue ?? '-'}`
                      : (field.enumValues ?? []).join(', ') || '—'}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          <Box component="form" onSubmit={onCreateField} sx={{ mt: 2 }}>
            <Stack direction={{ xs: 'column', md: 'row' }} spacing={1}>
              <TextField
                label="Field key"
                size="small"
                value={fieldForm.fieldKey}
                onChange={(e) => setFieldForm({ ...fieldForm, fieldKey: e.target.value })}
                required
              />
              <TextField
                label="Label"
                size="small"
                value={fieldForm.label}
                onChange={(e) => setFieldForm({ ...fieldForm, label: e.target.value })}
                required
              />
              <TextField
                select
                label="Type"
                size="small"
                value={fieldForm.dataType}
                onChange={(e) => setFieldForm({ ...fieldForm, dataType: e.target.value })}
                sx={{ minWidth: 120 }}
              >
                {DATA_TYPES.map((item) => (
                  <MenuItem key={item} value={item}>
                    {item}
                  </MenuItem>
                ))}
              </TextField>
              {fieldForm.dataType === 'ENUM' && (
                <TextField
                  label="Enum values (comma)"
                  size="small"
                  value={fieldForm.enumValuesText}
                  onChange={(e) => setFieldForm({ ...fieldForm, enumValuesText: e.target.value })}
                />
              )}
              <Button
                type="submit"
                variant="outlined"
                disabled={!currentVersion || createField.isPending}
              >
                Add field
              </Button>
            </Stack>
          </Box>
        </Paper>
      )}

      {activeId && (
        <Paper sx={{ p: 2, border: '1px solid #D5DEE8' }}>
          <Typography variant="h6" sx={{ mb: 1 }}>
            Validate sample values
          </Typography>
          <TextField
            fullWidth
            multiline
            minRows={5}
            value={validateJson}
            onChange={(e) => setValidateJson(e.target.value)}
            sx={{ fontFamily: 'IBM Plex Mono, monospace' }}
          />
          <Button sx={{ mt: 1 }} variant="contained" onClick={onValidate}>
            Validate
          </Button>
          {validateResult && (
            <Alert sx={{ mt: 2 }} severity={validateResult.valid ? 'success' : 'error'}>
              {validateResult.valid ? 'Valid' : validateResult.errors.join(' | ') || 'Invalid'}
            </Alert>
          )}
        </Paper>
      )}
    </Stack>
  )
}
