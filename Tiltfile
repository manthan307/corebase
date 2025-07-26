k8s_yaml([
    'k8s/auth.yml',
    'k8s/postgres.yml',
    'k8s/postgres-secret.yml',
    # 'k8s/storage.yml',
    # 'k8s/realtime.yml',
])


docker_build('auth', './auth',
 dockerfile='./auth/Dockerfile.dev',
 live_update=[
    sync('./auth', '/app'),
    run('npm install', trigger='./auth/package.json')
  ])
# docker_build('storage', './storage')
# docker_build('realtime', './realtime')

k8s_resource('postgres')
k8s_resource('auth', port_forwards=3000)
# k8s_resource('storage', port_forwards=4000)
# k8s_resource('realtime', port_forwards=5000)


# # For Go
# live_update(
#     'storage',
#     [
#         sync('./storage', '/app'),
#         run('go build -o app .'),
#     ]
# )
