import Box from '@mui/material/Box';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { observer } from 'mobx-react-lite';
import { configStore } from '../store/configuration-store';
import { PullRequest } from '../models/pull-request';

const columns: GridColDef<PullRequest>[] = [
  { field: 'number', headerName: 'ID', width: 90, type: 'number' },
  { field: 'title', headerName: 'Title', width: 400, editable: false },
  { field: 'url', headerName: 'URL', width: 499, editable: false },
  { field: 'state', headerName: 'Status', width: 110, editable: false },
  { field: 'createdAt', headerName: 'Created At', width: 230, editable: false },
];

// Reads the store during render, so a second fetch updates the grid. The
// previous version copied into state from an effect that ran only on mount.
const PullRequestDataGrid = observer(() => (
  <Box sx={{ height: 400, width: '100%' }}>
    <DataGrid
      getRowId={(row) => row.number}
      rows={configStore.pullRequests}
      columns={columns}
      initialState={{ pagination: { paginationModel: { pageSize: 5 } } }}
      pageSizeOptions={[5]}
      disableRowSelectionOnClick
    />
  </Box>
));

export default PullRequestDataGrid;
