from fastapi import FastAPI
from football.routers.hello_world import router as hello_world

apiv1 = FastAPI()
apiv1.include_router(hello_world)