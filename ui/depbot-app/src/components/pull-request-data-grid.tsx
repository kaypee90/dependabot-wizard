import { useEffect, useState } from 'react';
import Box from '@mui/material/Box';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { configStore } from '../store/configuration-store';
import { PullRequest } from '../models/pull-request';


const PullRequestDataGrid = () => {
  const [rows, setRows] = useState<PullRequest[]>([]);

  const columns: GridColDef<(typeof rows)[number]>[] = [
    {
      field: 'number',
      headerName: 'ID',
      width: 90,
      type: 'number'
    },
    {
      field: 'title',
      headerName: 'Title',
      width: 400,
      editable: false,
    },
    {
      field: 'url',
      headerName: 'URL',
      width: 499,
      editable: false,
    },
    {
      field: 'state',
      headerName: 'Status',
      width: 110,
      editable: false,
    },
    {
      field: 'createdAt',
      headerName: 'Created At',
      width: 230,
      editable: false,
    },
  ];


  useEffect(() => {
    setRows(configStore.pullRequests);
  }, []);

  return (
    <Box sx={{ height: 400, width: '100%' }}>
      <DataGrid
        getRowId={(row) => row.number}
        rows={rows}
        columns={columns}
        initialState={{
          pagination: {
            paginationModel: {
              pageSize: 5,
            },
          },
        }}
        pageSizeOptions={[5]}
        disableRowSelectionOnClick
      />
    </Box>
  );
}

export default PullRequestDataGrid;
