import React from "react";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import CheckIcon from '@mui/icons-material/Check';
import { observer } from "mobx-react-lite";

import { configStore } from '../store/configuration-store';
import { palette, MONO } from '../utils/theme';

const EmptyState = () => (
  <Box sx={{ p: 4, textAlign: 'center' }}>
    <Typography variant="body2" sx={{ color: palette.textFaint, maxWidth: 380, mx: 'auto' }}>
      Scan the project or choose an ecosystem, a directory and a schedule. The file appears here as
      you go, and is checked before it can be written.
    </Typography>
  </Box>
);

/**
 * The generated file, presented as a code panel rather than loose text: a
 * header naming the path it will be written to, line numbers, and a copy
 * action that works whether or not a depbot process is behind the page.
 */
const ConfigYaml: React.FC = observer(() => {
  const [copied, setCopied] = React.useState(false);
  const yaml = configStore.yamlData;
  const lines = yaml ? yaml.replace(/\n$/, '').split('\n') : [];

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(yaml);
      setCopied(true);
      setTimeout(() => setCopied(false), 1800);
    } catch {
      setCopied(false);
    }
  };

  return (
    <Box sx={{ border: `1px solid ${palette.border}`, borderRadius: 2, overflow: 'hidden', background: palette.surface }}>
      <Stack
        direction="row"
        alignItems="center"
        spacing={1}
        sx={{
          px: 2, py: 1.25,
          borderBottom: `1px solid ${palette.borderSoft}`,
          background: palette.raised,
        }}
      >
        <Typography
          sx={{ flexGrow: 1, fontFamily: MONO, fontSize: '0.78rem', color: palette.textMuted }}
        >
          .github/dependabot.yml
        </Typography>
        {yaml && (
          <>
            <Typography variant="caption" sx={{ color: palette.textFaint }}>
              {lines.length} lines
            </Typography>
            <Tooltip title={copied ? 'Copied' : 'Copy'}>
              <IconButton size="small" onClick={copy} aria-label="Copy configuration">
                {copied
                  ? <CheckIcon sx={{ fontSize: 17, color: palette.success }} />
                  : <ContentCopyIcon sx={{ fontSize: 16, color: palette.textMuted }} />}
              </IconButton>
            </Tooltip>
          </>
        )}
      </Stack>

      <Box sx={{ maxHeight: { md: 'calc(100vh - 190px)' }, overflow: 'auto' }}>
        {lines.length === 0 ? <EmptyState /> : (
          <Box component="pre" sx={{ m: 0, p: 2, fontSize: '0.8rem', lineHeight: 1.7 }}>
            {lines.map((line, index) => (
              <Box component="div" key={index} sx={{ display: 'flex' }}>
                <Box
                  component="span"
                  sx={{
                    width: 30, flexShrink: 0, textAlign: 'right', pr: 2,
                    color: palette.textFaint, userSelect: 'none',
                  }}
                >
                  {index + 1}
                </Box>
                <Box component="span" sx={{ color: palette.text, whiteSpace: 'pre' }}>
                  {line || ' '}
                </Box>
              </Box>
            ))}
          </Box>
        )}
      </Box>
    </Box>
  );
});

export default ConfigYaml;
