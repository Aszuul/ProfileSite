FROM golang:1.25 AS wasm-builder
WORKDIR /src
COPY go ./go
COPY build_wasm.sh ./
RUN chmod +x build_wasm.sh && ./build_wasm.sh

FROM python:3.13-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
COPY --from=wasm-builder /src/static/wasm/ /app/static/wasm/
EXPOSE 8080
CMD ["python", "app.py"]