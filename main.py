import uvicorn
from football.app import create_app

app = create_app()

def main():
    config = uvicorn.Config("main:app", port=8000, log_level="info")
    server = uvicorn.Server(config)
    server.run()


if __name__ == '__main__':
    main()
