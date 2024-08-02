export interface PlanDetail {
  description: string;
}

export interface PlanTemplate {
  details: PlanDetail[];
}

const getPlanTemplate = (): PlanTemplate => {
  return {
    details: [
      {
        description: "{{includedCredits}} compute tokens per month included (about {{numMolecules}} molecules)"
      },
      {
        description: "Every additional token costs {{overageCharge}} USD"
      },
      {
        description: "Cancel subscription any time"
      }
    ]
  };
};

export default getPlanTemplate;