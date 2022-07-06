.PHONY: run requirements

run:
	gunicorn -w 4 -k uvicorn.workers.UvicornWorker main:app

requirements:
	poetry export --dev -f requirements.txt --output requirements.txt