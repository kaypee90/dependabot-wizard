import React from "react";
import { Box, Typography } from "@mui/material";


type Props = {
  code: string
}

export default function ConfigYaml({ code }: Props) {
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
      background: '#fff'
    }}>
      {code}
    </Box>
  );
};
