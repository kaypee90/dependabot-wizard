import React from "react";
import { Box } from "@mui/material";
import { observer } from "mobx-react-lite";
import { configStore } from '../store/configuration-store';

const ConfigYaml: React.FC = observer(() => {
  return (
    <Box sx={{
      whiteSpace: 'pre',
      textShadow: 'white 0px 1px',
      fontFamily: 'Consolas, Monaco, "Andale Mono", "Ubuntu Mono", monospace',
      fontSize: '1em',
      textAlign: 'left',
      wordSpacing: 'normal',
      wordBreak: 'normal',
      lineHeight: '1.5',
      padding: '1em',
      tabSize: '4',
      hyphens: 'none',
      overflow: 'auto',
      margin: '0px 0px',
    }}>
      {configStore.yamlData}
    </Box>
  )
});

export default ConfigYaml;
