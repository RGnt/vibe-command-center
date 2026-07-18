FROM postgres:16

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    git \
    postgresql-server-dev-16 \
    flex \
    bison \
    && rm -rf /var/lib/apt/lists/*

# Clone and build Apache AGE
# AGE v1.5.0 supports PostgreSQL 16
RUN git clone -b PG16/v1.5.0-rc0 https://github.com/apache/age.git /age \
    && cd /age \
    && make \
    && make install \
    && cd / \
    && rm -rf /age

# Clone and build pgvector
# Using latest release tag (v0.7.2 or just main, we'll use a stable tag if possible, or just the main branch)
RUN git clone --branch v0.7.0 https://github.com/pgvector/pgvector.git /pgvector \
    && cd /pgvector \
    && make \
    && make install \
    && cd / \
    && rm -rf /pgvector

# Add extensions to shared_preload_libraries
RUN echo "shared_preload_libraries = 'age,vector'" >> /usr/share/postgresql/postgresql.conf.sample

# Set default command
CMD ["postgres"]
