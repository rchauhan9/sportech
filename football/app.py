from fastapi import FastAPI
from football.routers.subapp import apiv1


def create_app() -> FastAPI:
    app = FastAPI()
    app.mount("/api/v1", apiv1)
    return app
