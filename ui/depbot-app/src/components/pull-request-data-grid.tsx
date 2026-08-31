import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import Chip from '@mui/material/Chip';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { observer } from 'mobx-react-lite';

import { configStore } from '../store/configuration-store';
import { PullRequest } from '../models/pull-request';
import { palette } from '../utils/theme';

const STATE_COLOUR: Record<string, string> = {
  OPEN: palette.accent,
  MERGED: '#A78BD9',
  CLOSED: palette.textFaint,
};

const columns: GridColDef<PullRequest>[] = [
  {
    field: 'number',
    headerName: '#',
    width: 80,
    type: 'number',
    renderCell: (params) => (
      <Link href={params.row.url} target="_blank" rel="noreferrer"
        sx={{ color: palette.accent, textDecoration: 'none' }}>
        {params.value}
      </Link>
    ),
  },
  { field: 'title', headerName: 'Title', flex: 1, minWidth: 280 },
  {
    field: 'state',
    headerName: 'Status',
    width: 120,
    renderCell: (params) => (
      <Chip
        size="small"
        label={String(params.value ?? '').toLowerCase()}
        sx={{
          background: 'transparent',
          color: STATE_COLOUR[params.value as string] ?? palette.textMuted,
          border: `1px solid ${STATE_COLOUR[params.value as string] ?? palette.border}`,
        }}
      />
    ),
  },
  {
    field: 'createdAt',
    headerName: 'Opened',
    width: 130,
    valueFormatter: (value) => (value ? String(value).split('T')[0] : ''),
  },
];

// Reads the store during render, so a second fetch updates the grid. The
// previous version copied into state from an effect that ran only on mount.
const PullRequestDataGrid = observer(() => (
  <Box sx={{ height: 420, width: '100%' }}>
    <DataGrid
      getRowId={(row) => row.number}
      rows={configStore.pullRequests}
      columns={columns}
      initialState={{ pagination: { paginationModel: { pageSize: 10 } } }}
      pageSizeOptions={[10, 25]}
      disableRowSelectionOnClick
      sx={{
        border: `1px solid ${palette.border}`,
        borderRadius: 2,
        background: palette.surface,
        '& .MuiDataGrid-columnHeaders': { background: palette.raised },
        '& .MuiDataGrid-cell': { borderColor: palette.borderSoft },
        '& .MuiDataGrid-footerContainer': { borderColor: palette.borderSoft },
      }}
    />
  </Box>
));

export default PullRequestDataGrid;
