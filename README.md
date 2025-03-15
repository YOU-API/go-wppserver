The easiest way to run this project is to download the [pre-built binaries](https://github.com/YOU-API/go-wppserver/releases) in `bin.zip` file. After that, follow the instructions in the **readme.txt** file in the unzipped folder:

**Step 1**

Set up the `.env` file stored in the `\bin` directory or leave it as it is.

**Step 2**

Locate the corresponding binary file in the folder. For example, on my machine, I use `wppserver-windows-amd64.exe`.

**Step 3**

Run the chosen binary.

**Step 4**

Test all API requests using the Swagger interface at `http://localhost:8786/docs/api/`. The initial authentication user is defined in the `.env` file.

Note: The host and port for access should match the ones defined in the `.env` file.

**Step 5**

Through the `/device/login` route get a base64 image of the qrcode. Paste the base64 string into a web browser and launch the qrcode with your whatsapp.

**Step 6**

After connecting your device to the API, try sending a text message with the `/chat/send/text` route.
