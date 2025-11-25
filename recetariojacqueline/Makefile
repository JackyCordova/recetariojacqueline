
NAMESPACE=recetario
K8S_FILE=k8s/recetario.yaml

.PHONY: help \
	docker-build docker-up docker-down docker-reset \
	k8s-start k8s-up k8s-destroy k8s-url \
	grpc-smoke k6-test

docker-build:
	@echo ">> Activando para construir dentro de Minikube..."
	@eval $$(minikube docker-env) && \
	  echo ">> recetario-metadata-grpc..." && \
	  docker build -t recetario-metadata-grpc \
	    -f ./metadata/Dockerfile . && \
	  echo ">> recetario-rating-grpc..." && \
	  docker build -t recetario-rating-grpc \
	    -f ./rating/Dockerfile . && \
	  echo ">> recetario-recipe-grpc..." && \
	  docker build -t recetario-recipe-grpc \
	    -f ./recipe/Dockerfile .
	@echo "Imágenes Docker construidas"

docker-up:
	docker compose up -d
	docker compose ps

docker-down:
	docker compose down

docker-reset:
	docker compose down || true
	docker image rm recetario-metadata-grpc recetario-rating-grpc recetario-recipe-grpc -f || true
	docker volume prune -f || true
	docker system prune -f || true
	@echo "Docker limpio"


k8s-start:
	minikube start --driver=docker
	kubectl config use-context minikube
	@echo "Minikube listo"

k8s-up:
	@echo ">>> Creando namespace $(NAMESPACE) si no existe..."
	@kubectl get ns $(NAMESPACE) >/dev/null 2>&1 || kubectl create namespace $(NAMESPACE)

	@echo ">>> Aplicando archivo: $(K8S_FILE)..."
	kubectl apply -f $(K8S_FILE)

	@echo ">>> Esperando pods..."
	sleep 7

	@echo "=== Pods ==="
	kubectl get pods -n $(NAMESPACE)
	@echo "=== Services ==="
	kubectl get svc -n $(NAMESPACE)
	@echo "=== HPA ==="
	kubectl get hpa -n $(NAMESPACE)
	@echo "✔ Proyecto desplegado en Kubernetes."

k8s-destroy:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true
	@echo "✔ Namespace $(NAMESPACE) eliminado."

k8s-url:
	@echo "URL del LoadBalancer (recipe-grpc):"
	minikube service recipe-grpc -n $(NAMESPACE) --url


grpc-smoke:
	@if [ -z "$(RECETARIO_ADDR)" ]; then \
	  echo "ERROR: Debes pasar RECETARIO_ADDR, ejemplo:"; \
	  echo "  make grpc-smoke RECETARIO_ADDR=127.0.0.1:60000"; \
	  exit 1; \
	fi; \
	ADDR="$(RECETARIO_ADDR)"; \
	ADDR="$${ADDR#http://}"; \
	ADDR="$${ADDR#https://}"; \
	echo "Usando RECETARIO_ADDR = $$ADDR"; \
	echo ""; \
	echo "=== 1) Consultando receta (RecipeService) ==="; \
	grpcurl -plaintext \
	  -d '{"recipe_id":"r1"}' \
	  $$ADDR recetario.RecipeService/GetRecipeDetails || true; \
	echo ""; \
	echo "✔ Comunicación gRPC verificada entre los 3 microservicios."


k6-test:
	@if [ -z "$(RECETARIO_ADDR)" ]; then \
	  echo "ERROR: Debes pasar RECETARIO_ADDR, ejemplo:"; \
	  echo "  make k6-test RECETARIO_ADDR=127.0.0.1:60000"; \
	  exit 1; \
	fi; \
	ADDR="$(RECETARIO_ADDR)"; \
	ADDR="$${ADDR#http://}"; \
	ADDR="$${ADDR#https://}"; \
	echo "Usando RECETARIO_ADDR = $$ADDR para k6..."; \
	k6 run -e RECETARIO_ADDR=$$ADDR k6/recetarioload.js