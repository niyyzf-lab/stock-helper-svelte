<script lang="ts">
  import { fade, fly } from 'svelte/transition'
  import { quintOut } from 'svelte/easing'
  import { CakeIcon } from 'lucide-svelte'
  import StockChart from '../StockChart.svelte'
  import { toastStore } from '../../stores/toast'
  import type { TestResult, StockQuiz } from 'src/types/stock'

  export let currentStock: any
  export let currentStockIndex: number
  export let stocks: any[]
  export let showActions: boolean
  export let difficulty: number
  export let currentQuiz: StockQuiz | null

  let showQuizModal = false
  let currentQuestionIndex = 0
  let userAnswers: { [key: number]: string } = {}
  let showExplanation = false
  let quizScore = 0

  $: progress = ((currentStockIndex + 1) / stocks.length * 100).toFixed(0)

  // 重置状态
  function resetQuizState() {
    showActions = false
    console.log('状态已重置', { showActions })
  }

  // 监听状态变化
  $: {
    console.log('状态变化', { currentQuiz, showActions })
  }

  // 格式化日期
  function formatDate(dateStr: string, timeStr?: string) {
    if (!dateStr) return '未知时间'
    try {
      const date = new Date(dateStr)
      let result = date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
      if (timeStr) {
        result += ' ' + timeStr
      }
      return result
    } catch {
      return '无效日期'
    }
  }

  function handleAnswer(answer: string) {
    if (showExplanation) return
    
    const currentQuestion = currentQuiz?.quizzes[currentQuestionIndex]
    if (!currentQuestion) return
    
    userAnswers[currentQuestionIndex] = answer
    showExplanation = true
    
    // 计算得分
    if (answer === currentQuestion.correctAnswer) {
      quizScore += Math.floor(100 / currentQuiz!.quizzes.length)
    }
  }
  
  function nextQuestion() {
    if (!currentQuiz) return
    
    if (currentQuestionIndex < currentQuiz.quizzes.length - 1) {
      currentQuestionIndex++
      showExplanation = false
    } else {
      // 答题完成，显示最终得分
      showQuizModal = false
      toastStore.success(`测试完成！得分：${quizScore}分`)
    }
  }
  
  function startQuiz() {
    currentQuestionIndex = 0
    userAnswers = {}
    showExplanation = false
    quizScore = 0
    showQuizModal = true
  }
</script>

<div class="test-container">
  <div class="chart-container">
    <StockChart
      code={currentStock.dm}
      endTime={currentStock.testDate}
      freq="dh"
    />
    <!-- 答题按钮只在有题目时显示 -->
    {#if currentQuiz && showActions}
      <button 
        class="quiz-button"
        on:click={startQuiz}
        in:fly|local={{y: -20, duration: 300, delay: 200}}
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="quiz-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3"/>
          <circle cx="12" cy="12" r="1"/>
          <circle cx="7" cy="12" r="1"/>
          <circle cx="17" cy="12" r="1"/>
        </svg>
        <span class="quiz-button-text">开始答题</span>
      </button>
    {/if}
  </div>

  <div class="test-footer">
    <div class="footer-content">
      <!-- 左侧股票信息 -->
      <div class="stock-info">
        <div class="stock-basic">
          <span class="stock-code">{currentStock.dm}</span>
          <span class="divider">·</span>
          <span class="stock-name">{currentStock.mc}</span>
        </div>
        <span class="test-time">
          {formatDate(currentStock.testDate)}
        </span>
      </div>

      <!-- 中间显示AI出题状态 -->
      <div class="quiz-status">
        {#if !currentQuiz}
          <button 
            class="retry-button"
            on:click={() => showActions = true}
          >
            重试获取题目
          </button>
        {/if}
      </div>

      <!-- 右侧进度信息 -->
      <div class="progress-info">
        <CakeIcon size={14} strokeWidth={1.5} />
        <span class="progress-text">
          <span class="current">{currentStockIndex + 1}</span>
          <span class="divider">/</span>
          <span class="total">{stocks.length}</span>
        </span>
      </div>
    </div>
  </div>
</div>

<!-- 答题模态框 -->
{#if showQuizModal && currentQuiz}
  <div class="quiz-modal" transition:fade>
    <div class="quiz-modal-content">
      <div class="quiz-header">
        <h3>第 {currentQuestionIndex + 1}/{currentQuiz.quizzes.length} 题</h3>
        <button class="close-button" on:click={() => showQuizModal = false}>×</button>
      </div>
      
      {#if currentQuiz.quizzes[currentQuestionIndex]}
        {@const question = currentQuiz.quizzes[currentQuestionIndex]}
        <div class="quiz-body">
          <h4 class="quiz-title">{question.title}</h4>
          <p class="quiz-question">{question.question}</p>
          
          <div class="quiz-options">
            {#each question.options as option}
              <button
                class="option-button"
                class:selected={userAnswers[currentQuestionIndex] === option.id}
                class:correct={showExplanation && option.id === question.correctAnswer}
                class:incorrect={showExplanation && userAnswers[currentQuestionIndex] === option.id && option.id !== question.correctAnswer}
                on:click={() => handleAnswer(option.id)}
                disabled={showExplanation}
              >
                <span class="option-id">{option.id}</span>
                <span class="option-content">{option.content}</span>
              </button>
            {/each}
          </div>
          
          {#if showExplanation}
            <div class="explanation" in:fly={{y: 20, duration: 300}}>
              <h5>解释</h5>
              <p>{question.explanation}</p>
              <button class="next-button" on:click={nextQuestion}>
                {currentQuestionIndex < currentQuiz.quizzes.length - 1 ? '下一题' : '完成'}
              </button>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .test-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    gap: 1px;
    background: var(--neutral-100);
    border-radius: var(--radius-lg);
    overflow: hidden;
  }

  .chart-container {
    flex: 1;
    min-height: 0;
    position: relative;
    background: white;
    overflow: hidden;
  }

  .test-footer {
    background: white;
    padding: 12px 16px;
    border-top: 1px solid var(--border-color);
    border-bottom-left-radius: var(--radius-lg);
    border-bottom-right-radius: var(--radius-lg);
  }

  .footer-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .stock-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .stock-basic {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .stock-code {
    font-family: var(--font-mono);
    font-size: 15px;
    font-weight: 600;
    color: var(--primary-600);
  }

  .stock-name {
    font-size: 15px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .test-time {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .divider {
    color: var(--neutral-300);
  }

  .progress-info {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    background: var(--neutral-50);
    border: 1px solid var(--border-color);
    border-radius: 24px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .progress-text {
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .progress-info .current {
    color: var(--primary-600);
    font-weight: 600;
  }

  .progress-info .total {
    color: var(--text-secondary);
  }

  .quiz-button {
    position: absolute;
    top: 20px;
    right: 20px;
    padding: 10px 20px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 20px;
    font-weight: 500;
    font-size: 14px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition: all 0.2s;
    opacity: 0.9;
    backdrop-filter: blur(4px);
    z-index: 101;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .quiz-button:hover {
    background: #2563eb;
    transform: translateY(-2px);
    opacity: 1;
    box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
  }

  .quiz-icon {
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0% { transform: scale(1); }
    50% { transform: scale(1.1); }
    100% { transform: scale(1); }
  }

  .quiz-button-text {
    display: inline-block;
    font-weight: 600;
    letter-spacing: 0.5px;
  }

  .quiz-status {
    display: flex;
    align-items: center;
    gap: 8px;
    margin: 0 16px;
  }

  .retry-button {
    padding: 6px 12px;
    border-radius: 16px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .quiz-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .quiz-modal-content {
    background: white;
    padding: 0;
    border-radius: 12px;
    min-width: 600px;
    max-width: 90vw;
    max-height: 90vh;
    overflow-y: auto;
  }

  .quiz-header {
    padding: 16px 24px;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .quiz-body {
    padding: 24px;
  }

  .quiz-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 12px;
  }

  .quiz-question {
    color: var(--text-secondary);
    margin-bottom: 24px;
    line-height: 1.6;
  }

  .quiz-options {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .option-button {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    background: white;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
    position: relative;
    overflow: hidden;
  }

  .option-button:hover:not(:disabled) {
    background: var(--neutral-50);
    border-color: var(--primary-400);
  }

  .option-button.selected {
    background: #f0f9ff;
    border-color: #3b82f6;
  }

  .option-button.correct {
    background: #f0fdf4;
    border-color: #22c55e;
    color: #166534;
  }

  .option-button.incorrect {
    background: #fef2f2;
    border-color: #ef4444;
    color: #991b1b;
  }

  .option-button.correct::after,
  .option-button.incorrect::after {
    content: "";
    position: absolute;
    right: 16px;
    top: 50%;
    transform: translateY(-50%);
    width: 20px;
    height: 20px;
    background-repeat: no-repeat;
    background-position: center;
    background-size: contain;
  }

  .option-button.correct::after {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2322c55e'%3E%3Cpath d='M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z'/%3E%3C/svg%3E");
  }

  .option-button.incorrect::after {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%23ef4444'%3E%3Cpath d='M24 20.188l-8.315-8.209 8.2-8.282-3.697-3.697-8.212 8.318-8.31-8.203-3.666 3.666 8.321 8.24-8.206 8.313 3.666 3.666 8.237-8.318 8.285 8.203z'/%3E%3C/svg%3E");
  }

  .option-button.correct.selected {
    background: #dcfce7;
  }

  .option-button.incorrect.selected {
    background: #fee2e2;
  }

  .option-id {
    font-weight: 600;
    color: inherit;
    min-width: 24px;
    opacity: 0.8;
  }

  .option-content {
    flex: 1;
  }

  .explanation {
    margin-top: 24px;
    padding: 16px;
    background: #f8fafc;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
  }

  .explanation h5 {
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: 8px;
    color: #1e293b;
  }

  .explanation p {
    color: #475569;
    line-height: 1.6;
  }

  .next-button {
    margin-top: 16px;
    padding: 8px 24px;
    background: var(--primary-500);
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .next-button:hover {
    background: var(--primary-600);
  }

  .close-button {
    font-size: 24px;
    color: var(--text-secondary);
    background: none;
    border: none;
    cursor: pointer;
    padding: 4px;
    line-height: 1;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style> 