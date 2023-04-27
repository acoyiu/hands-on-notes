FROM node:14.7.0

# 建立一個目錄 /app
WORKDIR /app 

# run 用於準備環境的命令
RUN cd /app
RUN npm install -g http-server

# 將本機端與 Dockerfile 同一層的所有檔案加到 Linux 的 /app 目錄底下
ADD ./dist /app

# # 打開 port
# EXPOSE 8080

# 正式命令
CMD http-server -c-1 -p 8080

### command for build =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
# docker build --tag sampleimage:v0.3 --no-cache .
# docker run -dit --name asd -p 8080:8080 sampleimage:v0.3

### command for export image tar file =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
# docker save -o image.tar sampleimage:v0.3

### command for import image tar file =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
# docker load -i ./image.tar

### add tag to image and update repository =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
# docker tag 59f3e3615488 sampleimage/myimage
# docker push sampleimage/myimage