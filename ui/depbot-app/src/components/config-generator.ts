import { configStore } from "../store/configuration-store";

export const generateYaml = (ecosystem: string, directory: string, schedule: string, reviewer: string, openPrLimit: number, label: string) => {
  const yamlData: string = `
  version: 2
  updates:
  - package-ecosystem: ${ecosystem}
    directory: ${directory}
    schedule:
      interval: ${schedule} 
    reviewers:
    - ${reviewer}
    open-pull-requests-limit: ${openPrLimit}
    labels:
    - ${label}`;

  configStore.saveYamlData(yamlData);
}
