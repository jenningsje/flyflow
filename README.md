# Flyflow 

> Flyflow is API middleware to optimize LLM applications, same response quality, 5x lower latency, secure, and much higher token limits

When building on top of LLMs, builders care about the following: 

- Response quality 
- Latency (both time to first token and tokens / second)
- Rate limits 
- Reliability
- Enterprise grade security

Flyflow is middleware designed to optimize for all of these qualities, built to be open source, high performance written in golang, and optionally self-hosted for maximum flexiblity.

## Fine tuning 

The flyflow completions API is a drop in replacement for the openai completions API. Use flyflow directly in your openai provider and start using the API. All openai features will work including embeddings. 

```python 
from openai import OpenAI

client = OpenAI(
 base_url="https://api.flyflow.dev",
 api_key='demo'
)

chat_completion = client.chat.completions.create(...)
```

Flyflow automatically tracks your query patterns with openai and you can use it to fine tune mixtral MoE or llama 70b to match the quality of GPT4 on your query patterns. 

## Inference 

Flyflow allows for substantially higher token limits and reliability by load balancing across many different inference providers. 

We host your custom fine turned models with providers like anyscale, together.ai, and fal, and optimize for latency, tokens / second, and rate limits, with a model that's the same level of quality as GPT4 for your queries. 

This also enables significantly higher reliability, because if a provider fails we can drop in fallbacks to pick up the load. 

## Security and observability

Flyflow can also act as security middleware, preventing sensitive information from reaching the inference provider (including openai and microsoft). 

We provide easy to configure plugins that allow you to filter PII from your queries, and advanced observability tools that help you understand how LLMs are being used by your organization.

# Configurability

Flyflow is designed to be extremely configurable. Want to just use us as security middleware, but run all of your inference through GPT4? We got you. Want to fine tune and back up you backend with 5 inference providers for the highest possible rate limits and tokens / second? No problem. 

Written in golang, our backend is designed to maximize for performance, without compromising on flexibility for developers. 


