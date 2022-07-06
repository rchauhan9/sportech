from fastapi import APIRouter

router = APIRouter(prefix="/hello", tags=["hello"], responses={404: {"description": "Not Found"}})


@router.get("/world")
async def world():
    return {"message": "Hello World"}
