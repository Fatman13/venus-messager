FROM filvenus/venus-buildenv AS buildenv

COPY . ./venus-messager
RUN export GOPROXY=https://goproxy.cn && cd venus-messager && make

RUN cd venus-messager && ldd ./venus-messager


FROM filvenus/venus-runtime

# DIR for app
WORKDIR /app

# copy the app from build env
COPY --from=buildenv  /go/venus-messager/venus-messager /app/venus-messager


# 拷贝依赖库
COPY --from=buildenv   /lib/x86_64-linux-gnu/libpthread.so.0 \
    /lib/x86_64-linux-gnu/libdl.so.2 \
    /lib/x86_64-linux-gnu/libc.so.6 \
    /lib/

COPY ./docker/script  /script

EXPOSE 39812

ENTRYPOINT ["/app/venus-messager","run"]
