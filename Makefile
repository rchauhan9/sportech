.PHONY: run requirements

run:
	poetry run python main.py

requirements:
	poetry export -f requirements.txt --output requirements.txt