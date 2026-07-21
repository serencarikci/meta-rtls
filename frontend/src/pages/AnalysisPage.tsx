import { FormEvent, useState } from 'react'
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

type Requirement = {
  tenantCode: string
  profileScale: string
  title: string
  priority: string
  expectedTags: number
  expectedEps: number
  retentionDays: number
}

type TenantProfile = {
  code: string
  name: string
  profileScale: string
  featureCodes: string[]
  metadataFields: number
  expectedTags: number
  expectedEps: number
  retentionDays: number
  requirementTitle: string
}

type CompareResponse = {
  profiles: TenantProfile[]
  notes: string[]
}

type ImpactResult = {
  title: string
  affectedTenants: number
  affectedEntities: number
  migrationRequired: boolean
  riskLevel: string
  complexityScore: number
  backwardCompatible: boolean
  summary: string
  savedId?: string
}

type ChangeRequest = {
  id: string
  title: string
  requestType: string
  riskLevel: string
  complexityScore: number
  migrationRequired: boolean
  status: string
}

const REQUEST_TYPES = ['ADD_METADATA_FIELD', 'REMOVE_METADATA_FIELD', 'CHANGE_FIELD_TYPE']

export default function AnalysisPage() {
  const token = useAuthStore((s) => s.token)
  const queryClient = useQueryClient()
  const [impactForm, setImpactForm] = useState({
    requestType: 'ADD_METADATA_FIELD',
    title: 'Add batteryHealth field',
    fieldKey: 'batteryHealth',
    dataType: 'NUMBER',
    isRequired: false,
    save: true,
  })
  const [impactResult, setImpactResult] = useState<ImpactResult | null>(null)

  const requirements = useQuery({
    queryKey: ['analysis-requirements'],
    queryFn: () => api<Requirement[]>('/api/v1/analysis/requirements', {}, token),
  })

  const compare = useQuery({
    queryKey: ['analysis-compare'],
    queryFn: () => api<CompareResponse>('/api/v1/analysis/compare', {}, token),
  })

  const changeRequests = useQuery({
    queryKey: ['analysis-change-requests'],
    queryFn: () => api<ChangeRequest[]>('/api/v1/analysis/change-requests', {}, token),
  })

  const runImpact = useMutation({
    mutationFn: () =>
      api<ImpactResult>(
        '/api/v1/analysis/impact',
        { method: 'POST', body: JSON.stringify(impactForm) },
        token,
      ),
    onSuccess: (data) => {
      setImpactResult(data)
      queryClient.invalidateQueries({ queryKey: ['analysis-change-requests'] })
    },
  })

  function onImpact(e: FormEvent) {
    e.preventDefault()
    runImpact.mutate()
  }

  return (
    <Stack spacing={3}>
      <div>
        <Typography variant="h4">Architecture analysis</Typography>
        <Typography color="text.secondary">
          Compare small / medium / large tenants and estimate metadata change cost.
        </Typography>
      </div>

      <Paper sx={{ p: 2 }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Requirement matrix
        </Typography>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Tenant</TableCell>
              <TableCell>Scale</TableCell>
              <TableCell>Title</TableCell>
              <TableCell>Priority</TableCell>
              <TableCell>Tags</TableCell>
              <TableCell>EPS</TableCell>
              <TableCell>Retention</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(requirements.data ?? []).map((item, index) => (
              <TableRow key={`${item.tenantCode}-${index}`}>
                <TableCell>{item.tenantCode}</TableCell>
                <TableCell>{item.profileScale}</TableCell>
                <TableCell>{item.title}</TableCell>
                <TableCell>{item.priority}</TableCell>
                <TableCell>{item.expectedTags}</TableCell>
                <TableCell>{item.expectedEps}</TableCell>
                <TableCell>{item.retentionDays}d</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Paper>

      <Paper sx={{ p: 2 }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Profile compare
        </Typography>
        <Stack direction={{ xs: 'column', md: 'row' }} gap={2}>
          {(compare.data?.profiles ?? []).map((profile) => (
            <Box
              key={profile.code}
              sx={{
                flex: 1,
                p: 2,
                border: '1px solid var(--line)',
                background: '#FCFAFD',
              }}
            >
              <Typography variant="overline">{profile.profileScale}</Typography>
              <Typography variant="h6">{profile.code}</Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                {profile.requirementTitle || profile.name}
              </Typography>
              <Typography variant="body2">Tags: {profile.expectedTags}</Typography>
              <Typography variant="body2">EPS: {profile.expectedEps}</Typography>
              <Typography variant="body2">Retention: {profile.retentionDays} days</Typography>
              <Typography variant="body2">Metadata fields: {profile.metadataFields}</Typography>
              <Stack direction="row" gap={0.5} flexWrap="wrap" sx={{ mt: 1 }}>
                {profile.featureCodes.map((code) => (
                  <Chip key={code} size="small" label={code} variant="outlined" />
                ))}
              </Stack>
            </Box>
          ))}
        </Stack>
        <Stack sx={{ mt: 2 }} spacing={0.5}>
          {(compare.data?.notes ?? []).map((note) => (
            <Typography key={note} variant="body2" color="text.secondary">
              • {note}
            </Typography>
          ))}
        </Stack>
      </Paper>

      <Paper sx={{ p: 2 }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Change impact analysis
        </Typography>
        <Box component="form" onSubmit={onImpact}>
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={1} sx={{ mb: 1 }}>
            <TextField
              select
              label="Type"
              size="small"
              value={impactForm.requestType}
              onChange={(e) => setImpactForm({ ...impactForm, requestType: e.target.value })}
              sx={{ minWidth: 220 }}
            >
              {REQUEST_TYPES.map((item) => (
                <MenuItem key={item} value={item}>
                  {item}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              label="Title"
              size="small"
              value={impactForm.title}
              onChange={(e) => setImpactForm({ ...impactForm, title: e.target.value })}
              required
              fullWidth
            />
            <TextField
              label="Field key"
              size="small"
              value={impactForm.fieldKey}
              onChange={(e) => setImpactForm({ ...impactForm, fieldKey: e.target.value })}
            />
            <Button type="submit" variant="contained" disabled={runImpact.isPending}>
              Analyze
            </Button>
          </Stack>
        </Box>

        {impactResult && (
          <Alert severity={impactResult.riskLevel === 'LOW' ? 'success' : 'warning'} sx={{ mt: 1 }}>
            {impactResult.summary}
            <br />
            Tenants: {impactResult.affectedTenants} · Entities: {impactResult.affectedEntities} ·
            Migration: {impactResult.migrationRequired ? 'yes' : 'no'} · Score:{' '}
            {impactResult.complexityScore}
          </Alert>
        )}
      </Paper>

      <Paper sx={{ p: 2 }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Saved change requests
        </Typography>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Title</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Risk</TableCell>
              <TableCell>Score</TableCell>
              <TableCell>Migration</TableCell>
              <TableCell>Status</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(changeRequests.data ?? []).map((item) => (
              <TableRow key={item.id}>
                <TableCell>{item.title}</TableCell>
                <TableCell>{item.requestType}</TableCell>
                <TableCell>{item.riskLevel}</TableCell>
                <TableCell>{item.complexityScore}</TableCell>
                <TableCell>{item.migrationRequired ? 'yes' : 'no'}</TableCell>
                <TableCell>{item.status}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Paper>
    </Stack>
  )
}
