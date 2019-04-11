import asyncio
import websockets
import cv2
import base64

async def sendimages(websocket):
    cap = cv2.VideoCapture(0)
    while True:
            retval, img = cap.read()
            retval, buffer = cv2.imencode('.jpg', img)
            img = base64.b64encode(buffer)
            #buffer = buffer.tostring()
            #print(buffer)
            await websocket.send(img)
            await asyncio.sleep(0.01)

async def recievecontrol(websocket):
    async for message in websocket:
        print(message)

async def hello():

    async with websockets.connect(
        'ws://localhost:3000/broadcaster') as websocket:
        
        send_task = asyncio.ensure_future(sendimages(websocket))
        recieve_task = asyncio.ensure_future(recievecontrol(websocket))

        done, pending = await asyncio.wait(
            [send_task, recieve_task],
            return_when=asyncio.FIRST_COMPLETED
        )
        for task in pending:
            task.cancel()

asyncio.get_event_loop().run_until_complete(hello())