import { makeAutoObservable } from "mobx";
import { PullRequest } from "../models/pull-request";
import { DetectedEcosystem } from "../models/detected-ecosystem";
import { GroupingStrategy } from "../components/config-generator";

export type ScanState = "idle" | "scanning" | "done" | "failed";

class ConfigurationStore {
  yamlData: string = "";

  ecosystem: string = "";
  directory: string = "";
  schedule: string = "";
  reviewer: string = "";
  openPrLimit: string = "";
  label: string = "";

  showGithubInfo: boolean = false;
  showGithubRepoDetails: boolean = false;
  hasSavedRepoDetails: boolean = false;

  pullRequests: PullRequest[] = [];
  /** What the server searched for, shown when nothing came back. */
  pullRequestFilter: string = "";

  /** Repository-wide noise controls, applied to every update block. */
  grouping: GroupingStrategy = "none";
  ignoreMajor: boolean = false;

  detected: DetectedEcosystem[] = [];
  useDetected: boolean = false;
  scanState: ScanState = "idle";
  scanError: string = "";

  constructor() {
    makeAutoObservable(this);
  }

  setPackageConfiguration(ecosystem: string, directory: string, schedule: string, reviewer: string, openPrLimit: string, label: string) {
    this.ecosystem = ecosystem;
    this.directory = directory;
    this.schedule = schedule;
    this.reviewer = reviewer;
    this.openPrLimit = openPrLimit;
    this.label = label;
  }

  /**
   * Stores the rendered config verbatim. It must not be re-indented or
   * trimmed here: doing so is what used to produce an unparseable file.
   */
  saveYamlData(newYamlData: string) {
    this.yamlData = newYamlData;
  }

  setShowGithubInfo(showInfo: boolean) {
    this.showGithubInfo = showInfo;
  }

  setShowGithubRepoDetails(showDetails: boolean) {
    this.showGithubRepoDetails = showDetails;
  }

  setHasSavedRepoDetails(hasSaved: boolean) {
    this.hasSavedRepoDetails = hasSaved;
  }

  setPullRequests(prs: PullRequest[], filter: string = "") {
    this.pullRequests = prs ?? [];
    this.pullRequestFilter = filter;
  }

  setScanState(state: ScanState, error: string = "") {
    this.scanState = state;
    this.scanError = error;
  }

  setDetected(detected: DetectedEcosystem[]) {
    this.detected = detected;
    this.useDetected = detected.length > 0;
  }

  setUseDetected(useDetected: boolean) {
    this.useDetected = useDetected;
  }

  setGrouping(grouping: GroupingStrategy) {
    this.grouping = grouping;
  }

  setIgnoreMajor(ignoreMajor: boolean) {
    this.ignoreMajor = ignoreMajor;
  }

  /** True when there is enough information to write a file. */
  get isValid(): boolean {
    if (this.useDetected) {
      return this.detected.length > 0 && Boolean(this.schedule);
    }

    return Boolean(this.ecosystem && this.directory && this.schedule);
  }
}

export const configStore = new ConfigurationStore();
