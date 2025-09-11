# AI Model Configuration Guide

## 🎯 **Dynamic Model Configuration**

Your bot now supports **dynamic AI model selection** via environment variables. You can choose the exact model for each AI service based on your needs (cost, performance, capabilities).

## 🤖 **Available AI Models**

### **Claude.ai Models**
Configure with `CLAUDE_MODEL` environment variable:

| Model | Speed | Cost | Capabilities | Best For |
|-------|--------|------|-------------|----------|
| `claude-3-haiku-20240307` | ⚡ Fastest | 💰 Cheapest | Good reasoning | **Default** - Quick analysis |
| `claude-3-sonnet-20240229` | 🚀 Fast | 💰💰 Mid | Better reasoning | Complex requirements |
| `claude-3-opus-20240229` | 🐌 Slowest | 💰💰💰 Most expensive | Best reasoning | Mission-critical analysis |

### **OpenAI Models**
Configure with `OPENAI_MODEL` environment variable:

| Model | Speed | Cost | Capabilities | Best For |
|-------|--------|------|-------------|----------|
| `gpt-3.5-turbo` | ⚡ Fast | 💰 Cheap | Good performance | **Default** - Cost-effective |
| `gpt-4` | 🚀 Medium | 💰💰💰 Expensive | Excellent reasoning | High-quality analysis |
| `gpt-4-turbo-preview` | 🚀 Medium | 💰💰 Mid-expensive | Latest improvements | Balanced performance |
| `gpt-4o` | ⚡ Fastest GPT-4 | 💰💰 Mid | Optimized GPT-4 | Speed + quality |

### **Google Gemini Models**
Configure with `GEMINI_MODEL` environment variable:

| Model | Speed | Cost | Capabilities | Best For |
|-------|--------|------|-------------|----------|
| `gemini-pro` | 🚀 Fast | 💰 Free | Good performance | **Default** - Free tier |
| `gemini-1.5-pro-latest` | 🚀 Fast | 💰💰 Paid | Latest improvements | Enhanced reasoning |
| `gemini-1.5-flash-latest` | ⚡ Fastest | 💰 Free/Low | Speed optimized | Quick responses |

## ⚙️ **Configuration Examples**

### **Cost-Optimized Setup (Recommended for Development)**
```bash
# Use fastest/cheapest models
CLAUDE_MODEL=claude-3-haiku-20240307
OPENAI_MODEL=gpt-3.5-turbo
GEMINI_MODEL=gemini-pro
```

### **Quality-Optimized Setup (Production/Enterprise)**
```bash
# Use best reasoning models
CLAUDE_MODEL=claude-3-opus-20240229
OPENAI_MODEL=gpt-4
GEMINI_MODEL=gemini-1.5-pro-latest
```

### **Balanced Setup (Recommended for Production)**
```bash
# Balance cost vs quality
CLAUDE_MODEL=claude-3-sonnet-20240229
OPENAI_MODEL=gpt-4o
GEMINI_MODEL=gemini-1.5-pro-latest
```

### **Speed-Optimized Setup**
```bash
# Prioritize fast responses
CLAUDE_MODEL=claude-3-haiku-20240307
OPENAI_MODEL=gpt-4o
GEMINI_MODEL=gemini-1.5-flash-latest
```

## 🔄 **Intelligent Fallback System**

The bot uses a **4-tier intelligent fallback** that tries models in this order:

1. **Claude** (Primary) - Uses your `CLAUDE_MODEL`
2. **OpenAI** (1st Fallback) - Uses your `OPENAI_MODEL`  
3. **Gemini** (2nd Fallback) - Uses your `GEMINI_MODEL`
4. **Rule-based** (Ultimate Fallback) - Always works

### **Smart Model Selection**
- **Claude**: Best for complex architectural decisions and detailed code analysis
- **OpenAI**: Most reliable and widely available, excellent for general software tasks
- **Gemini**: Good alternative with competitive performance, often free
- **Rule-based**: Guarantees the bot never fails, provides basic analysis

## 💡 **Model Selection Recommendations**

### **For Different Use Cases:**

#### **Startups/Small Teams (Cost-conscious)**
```bash
CLAUDE_MODEL=claude-3-haiku-20240307    # Fast & cheap
OPENAI_MODEL=gpt-3.5-turbo              # Most cost-effective
GEMINI_MODEL=gemini-pro                 # Free tier
```

#### **Enterprise/Production (Quality-focused)**
```bash
CLAUDE_MODEL=claude-3-opus-20240229     # Best reasoning
OPENAI_MODEL=gpt-4                      # Most reliable
GEMINI_MODEL=gemini-1.5-pro-latest     # Latest features
```

#### **High-Volume/SaaS (Balance)**
```bash
CLAUDE_MODEL=claude-3-sonnet-20240229   # Good balance
OPENAI_MODEL=gpt-4-turbo-preview        # Optimized cost/performance
GEMINI_MODEL=gemini-1.5-flash-latest   # Speed focused
```

## 🚀 **Setup Instructions**

### 1. **Choose Your Models** (based on use case above)

### 2. **Update Your `.env` File**
```bash
# Copy from .env.example and customize
cp .env.example .env

# Edit with your preferred models
CLAUDE_MODEL=your_chosen_claude_model
OPENAI_MODEL=your_chosen_openai_model
GEMINI_MODEL=your_chosen_gemini_model
```

### 3. **Test Different Models**
You can test different models by simply changing the environment variables and restarting the bot:
```bash
# Try different models
export CLAUDE_MODEL=claude-3-opus-20240229
export OPENAI_MODEL=gpt-4
./yordamchi-dev-bot
```

## 📊 **Model Performance Comparison**

| Task Type | Claude-3 Opus | GPT-4 | Claude-3 Sonnet | GPT-4-Turbo | Claude-3 Haiku | GPT-3.5-Turbo |
|-----------|---------------|-------|-----------------|-------------|----------------|----------------|
| **Code Analysis** | 🟢 Excellent | 🟢 Excellent | 🟡 Good | 🟢 Excellent | 🟡 Good | 🟡 Good |
| **Architecture Planning** | 🟢 Excellent | 🟢 Excellent | 🟢 Excellent | 🟢 Excellent | 🟡 Good | 🟡 Good |
| **Task Estimation** | 🟢 Excellent | 🟢 Excellent | 🟡 Good | 🟢 Excellent | 🟡 Good | 🟡 Good |
| **Speed** | 🔴 Slow | 🟡 Medium | 🟡 Medium | 🟢 Fast | 🟢 Very Fast | 🟢 Very Fast |
| **Cost** | 🔴 High | 🔴 High | 🟡 Medium | 🟡 Medium | 🟢 Low | 🟢 Very Low |

## 🔧 **Advanced Configuration**

### **Runtime Model Switching**
Models are loaded from environment variables at startup. To switch models:
1. Update your `.env` file
2. Restart the bot
3. The bot will use the new models immediately

### **Mixed Model Strategy**
You can use different quality tiers for the fallback chain:
```bash
# Best model for primary, cheaper for fallbacks
CLAUDE_MODEL=claude-3-opus-20240229     # Premium quality primary
OPENAI_MODEL=gpt-3.5-turbo              # Cost-effective fallback
GEMINI_MODEL=gemini-pro                 # Free backup
```

This gives you premium analysis when Claude is available, but keeps costs low for fallback scenarios.

## 🎯 **Summary**

✅ **Fully Configurable**: Choose exact models for each AI service  
✅ **Smart Defaults**: Works out-of-the-box with sensible defaults  
✅ **Cost Control**: Pick models based on your budget  
✅ **Quality Control**: Upgrade to premium models when needed  
✅ **Intelligent Fallback**: Never fails due to robust 4-tier system  
✅ **Runtime Flexibility**: Easy to change models by updating config