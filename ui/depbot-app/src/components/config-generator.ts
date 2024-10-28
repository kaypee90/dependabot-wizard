import { configStore } from "../store/configuration-store";

export const generateYaml = () => {
  const yamlData: string = `
  version: 2
  updates:
  - package-ecosystem: ${configStore.ecosystem}
    directory: ${configStore.directory}
    schedule:
      interval: ${configStore.schedule} 
    reviewers:
    - ${configStore.reviewer}
    open-pull-requests-limit: ${configStore.openPrLimit}
    labels:
    - ${configStore.label}`;

  configStore.saveYamlData(yamlData);
}
