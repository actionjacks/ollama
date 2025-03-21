#!/bin/sh
export OLLAMA_PORT=11436

if [ -z "$(which nvidia-smi)" ]; then
    echo "--- NVIDIA-SMI ---"
    echo "No access to GPU"
else
    echo "--- NVIDIA-SMI ---"
    echo "GPU detected"
    export CUDA_VISIBLE_DEVICES=0
fi

echo "Starting Ollama..."
export OLLAMA_HOST=0.0.0.0:$OLLAMA_PORT
ollama serve & # Start Ollama in the background

sleep 5 # Wait for it to initialize

echo "Downloading models..."

ollama pull llama3:8b | tee /var/log/ollama_pull.log

if [ -f "/app/context/doc.txt" ]; then
    echo "Creating custom model based on doc.txt..."
    ollama create netiskuba -f /app/context/doc.txt
fi

echo "Ollama is ready."
tail -f /dev/null  # Keep the container running