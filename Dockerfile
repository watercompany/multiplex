FROM alpine:3.4

WORKDIR /

COPY ./execs ./
COPY ./server .
COPY ./main .
COPY ./scripts/multiplex.sh .

RUN mkdir -p /output

RUN mkdir -p /mnt/ssd1/plotfiles/final
RUN mkdir -p /mnt/ssd1/plotfiles/temp
RUN mkdir -p /mnt/ssd1/plotfiles/temp2

RUN mkdir -p /mnt/ssd2/plotfiles/final
RUN mkdir -p /mnt/ssd2/plotfiles/temp
RUN mkdir -p /mnt/ssd2/plotfiles/temp2

RUN mkdir -p /mnt/ssd3/plotfiles/final
RUN mkdir -p /mnt/ssd3/plotfiles/temp
RUN mkdir -p /mnt/ssd3/plotfiles/temp2

RUN mkdir -p /mnt/ssd4/plotfiles/final
RUN mkdir -p /mnt/ssd4/plotfiles/temp
RUN mkdir -p /mnt/ssd4/plotfiles/temp2

RUN mkdir -p /mnt/ssd5/plotfiles/final
RUN mkdir -p /mnt/ssd5/plotfiles/temp
RUN mkdir -p /mnt/ssd5/plotfiles/temp2

RUN mkdir -p /mnt/ssd6/plotfiles/final
RUN mkdir -p /mnt/ssd6/plotfiles/temp
RUN mkdir -p /mnt/ssd6/plotfiles/temp2

RUN mkdir -p /mnt/ssd7/plotfiles/final
RUN mkdir -p /mnt/ssd7/plotfiles/temp
RUN mkdir -p /mnt/ssd7/plotfiles/temp2

RUN mkdir -p /mnt/ssd8/plotfiles/final
RUN mkdir -p /mnt/ssd8/plotfiles/temp
RUN mkdir -p /mnt/ssd8/plotfiles/temp2

RUN mkdir -p /media/cx/10tb

ENTRYPOINT ["sh","/multiplex.sh"]
