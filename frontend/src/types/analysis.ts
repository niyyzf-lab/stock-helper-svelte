// 分析进度状态
export interface AnalysisProgress {
  step: number;
  total: number;
  message: string;
  phase: 'data' | 'analysis' | '';
}

// 分析请求参数
export interface AnalysisRequest {
  code: string;
}

// 分析响应
export interface AnalysisResponse {
  type: 'quick' | 'full';  // 快速分析或完整分析
  data: StockAnalysis;
}

// 快速分析结果
export interface QuickAnalysisResult {
  action: string;
  timing: string;
  targetPrice: string;
  stopLoss: string;
  pe: string;
  pb: string;
  performanceScore: string;
  riskLevel: string;
  mainForce: {
    direction: string;
    strength: string;
  };
  recentTrend: string;
}

export interface StockAnalysis {
  fundamental: {
    companyProfile: {
      industryPosition: string;
      competitiveEdge: string[];
      businessModel: string;
      coreBusiness: string;
      growthStage: string;
      highlights: string[];
      concerns: string[];
    };
    financialAnalysis: {
      performanceScore: string;
      revenue: {
        growth: string;
        stability: string;
        quality: string;
        futureOutlook: string;
        comment: string;
      };
      profitability: {
        grossMargin: string;
        netMargin: string;
        roe: string;
        quality: string;
        comment: string;
      };
      financialHealth: {
        debtLevel: string;
        cashFlow: string;
        operatingStatus: string;
        healthScore: string;
        comment: string;
      };
    };
    valuationAnalysis: {
      currentValuation: {
        pe: string;
        pb: string;
        ps: string;
        valuationLevel: string;
        comment: string;
      };
      relativeValuation: {
        industryComparison: string;
        historyComparison: string;
        valueScore: string;
        comment: string;
      };
    };
  };
  trend: {
    priceTrend: {
      mainTrend: string;
      recentTrend: string;
      trendStrength: string;
      keyLevels: {
        support: string[];
        resistance: string[];
        majorLevel: string;
        comment: string;
      };
    };
    capitalFlow: {
      mainForce: {
        direction: string;
        strength: string;
        concentration: string;
        comment: string;
      };
      institutionalActivity: {
        trend: string;
        significantMove: string[];
        comment: string;
      };
    };
    marketSentiment: {
      shortTerm: string;
      mediumTerm: string;
      volatility: string;
      comment: string;
    };
  };
  risk: {
    coreRisks: {
      businessRisks: string[];
      financialRisks: string[];
      marketRisks: string[];
      riskLevel: string;
      comment: string;
    };
    specificRisks: {
      shortTerm: string[];
      mediumTerm: string[];
      longTerm: string[];
      comment: string;
    };
    riskMitigation: {
      suggestions: string[];
      keyMonitorings: string[];
      comment: string;
    };
  };
  trading: {
    newPosition: {
      action: string;
      timing: string;
      entryPoints: string[];
      positionSize: string;
      stopLoss: string;
      riskReward: string;
      comment: string;
    };
    existingPosition: {
      action: string;
      adjustmentType: string;
      exitPoints: string[];
      stopProfit: string;
      positionAdjust: string;
      comment: string;
    };
    investmentPlan: {
      timeHorizon: string;
      targetPrice: string;
      milestonePrice: string[];
      keyEvents: string[];
      comment: string;
    };
  };
} 