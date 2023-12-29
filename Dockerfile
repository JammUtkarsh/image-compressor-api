FROM golang:1.21 AS build
WORKDIR /app
COPY ./API . 
RUN go build -o apiService .


FROM golang:1.21 AS api-final
WORKDIR /app
COPY --from=build /app/apiService .
CMD ["./apiService"]
EXPOSE 8080


# Original File: https://github.com/cshum/imagor/blob/master/Dockerfile
ARG GOLANG_VERSION=1.21
FROM golang:${GOLANG_VERSION}-bookworm as img-final

ARG VIPS_VERSION=8.14.5
ARG TARGETARCH

WORKDIR /app

ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig

# Installs libvips + required libraries
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl \
  meson ninja-build pkg-config \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg62-turbo-dev libpng-dev \
  libwebp-dev libtiff-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev libopenjp2-7-dev liblcms2-dev \
  libgsf-1-dev libfftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libaom-dev \
  libheif-dev libspng-dev libcgif-dev && \
  cd /tmp && \
    curl -fsSLO https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.xz && \
    tar xf vips-${VIPS_VERSION}.tar.xz && \
    cd vips-${VIPS_VERSION} && \
    meson setup _build \
    --buildtype=release \
    --strip \
    --prefix=/usr/local \
    --libdir=lib \
    -Dgtk_doc=false \
    -Dmagick=disabled \
    -Dintrospection=false && \
    ninja -C _build && \
    ninja -C _build install && \
  ldconfig && \
  rm -rf /usr/local/lib/libvips-cpp.* && \
  rm -rf /usr/local/lib/*.a && \
  rm -rf /usr/local/lib/*.la

COPY ./ImageService .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=${TARGETARCH} go build -o goimg

CMD ["./goimg"]