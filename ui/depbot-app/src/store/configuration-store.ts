import { makeAutoObservable } from "mobx";
import { PullRequest } from "../models/pull-request";

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

  saveYamlData(newYamlData: string) {
    if (newYamlData.trim()) {
      this.yamlData = newYamlData.trim() + "\n";
    }
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

  setPullRequests(prs: PullRequest[]) {
    this.pullRequests = prs;
  }
}

export const configStore = new ConfigurationStore();

