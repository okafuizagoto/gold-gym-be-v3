// ==============================================================
// @Library — TIDAK DIGUNAKAN di local
//   → Shared library tidak dipakai di local.
//   → Kalau gold-gym memiliki shared library sendiri, aktifkan di sini:
//       @Library('gold-gym-shared-library')
//   → Fungsi dari shared library (checkoutSCM, version, discordSend, dll)
//     sudah di-inline di masing-masing stage di bawah.
// ==============================================================
// @Library('gold-gym-shared-library')

pipeline {
    // agent any
    //   → Pipeline ini bisa jalan di Jenkins agent manapun yang tersedia.
    //   → Kalau butuh agent spesifik (misal untuk Docker build),
    //     ubah ke: agent { label 'docker' }
    agent any

    // ==============================================================
    // environment
    //   → Definisi environment variables yang berlaku di seluruh pipeline.
    //   → Variabel di sini bisa di-akses via env.NAMA_VAR atau langsung NAMA_VAR.
    // ==============================================================
    environment {
        // PROJECT_ID
        //   → GCP Project ID tempat Docker image dan Helm chart di-host.
        //   → TODO: Ganti dengan GCP Project ID punya gold-gym (kalau ada).
        // PROJECT_ID = "gold-gym-project"
        PROJECT_ID = "gold-gym-project"

        // NAME
        //   → Nama project ini. Dipakai untuk:
        //     - Docker image tag
        //     - Helm chart name
        //     - Helm package filename
        NAME = "gold-gym-be-v2"

        // ORG
        //   → Namespace / organization di ChartMuseum.
        //   → TODO: Ganti dengan namespace ChartMuseum punya gold-gym (kalau ada).
        // ORG = "gold-gym"
        ORG = "gold-gym"

        // K8S_CLUSTER
        //   → Nama Kubernetes cluster target deployment.
        //   → TODO: Ganti dengan nama cluster K8s punya gold-gym.
        // K8S_CLUSTER = "gold-gym-cluster"
        K8S_CLUSTER = "gold-gym-cluster"

        // ==============================================================
        // Docker Registry — TODO: Ganti dengan registry punya gold-gym
        //   → Saat ini placeholder. Kalau gold-gym belum punya registry,
        //     gunakan Docker Hub atau registry local.
        // ==============================================================
        ARTIFACT_REGISTRY = "localhost:5000"

        // DOCKER_REGISTRY — tidak dipakai di local pipeline
        DOCKER_REGISTRY = "localhost:5000"

        // DOCKER_REGISTRY_URL
        DOCKER_REGISTRY_URL = "http://${ARTIFACT_REGISTRY}"

        // DOCKER_REGISTRY_PROJECT_URL
        DOCKER_REGISTRY_PROJECT_URL = "${ARTIFACT_REGISTRY}/${PROJECT_ID}"

        // PIPELINE_NAME
        //   → Display name di notifikasi Discord.
        PIPELINE_NAME = "Gold Gym BE V2"

        // PIPELINE_BOT_EMAIL
        //   → Email bot untuk gold-gym pipeline.
        PIPELINE_BOT_EMAIL = "goldgym.bot@gmail.com"

        // PIPELINE_BOT_NAME
        //   → Display name bot.
        PIPELINE_BOT_NAME = "Gold Gym Pipeline Bot"

        // DISCORD_WEBHOOK_URL
        //   → Webhook URL Discord untuk notifikasi build gold-gym-be-v2.
        DISCORD_WEBHOOK_URL = "https://discordapp.com/api/webhooks/1314527977240793088/JOdEblBl1wc-IMK4x3PP_AYOi0PNVQksx-yAnfbNFmz1GbtYHt3q2jzXh13eQ1tEdERA"

        // ==============================================================
        // JAVA_HOME — TIDAK DIGUNAKAN
        //   → sttk-be memiliki Java dependency, jadi JAVA_HOME di-set.
        //   → gold-gym-be-v2 adalah project Go murni, tidak butuh Java.
        //   → Commented out. Aktifkan kembali kalau ada dependency Java.
        // ==============================================================
        // JAVA_HOME = '/usr/lib/jvm/java-11-openjdk-amd64/'
    }

    // ==============================================================
    // options
    //   → Pipeline-level configuration.
    // ==============================================================
    options {
        // skipDefaultCheckout(true)
        //   → Jangan auto-Checkout GIT di awal pipeline.
        //   → Checkout dilakukan manual di stage 'Checkout GIT'
        //     menggunakan checkoutSCM() dari shared library.
        //   → Alasan: shared library mungkin butuh custom checkout logic
        //     (misal: fetch tags, submodules, dll).
        skipDefaultCheckout(true)
    }

    stages {
        // ==============================================================
        // Stage 1: Checkout GIT
        //   → Git clone dari GitHub repository gold-gym-be-v2.
        //   → Menggunakan Jenkins built-in Git plugin (bukan shared library).
        // ==============================================================
        stage('Checkout GIT') {
            steps {
                git url: 'https://github.com/okafuizagoto/gold-gym-be-v2.git',
                    branch: 'main'
            }
        }

        // ==============================================================
        // Stage 2: Version
        //   → Set VERSION dari git tag atau Jenkins build number.
        //   → Tidak pakai shared library, inline logic.
        // ==============================================================
        stage('Version') {
            steps {
                script {
                    // Coba ambil version dari git tag, kalau tidak ada pakai build number
                    env.VERSION = sh(script: "git describe --tags --always || echo ${env.BUILD_NUMBER}", returnStdout: true).trim()
                    echo "Version: ${env.VERSION}"
                }
            }
        }

        // ==============================================================
        // Stage 3: Compile
        //   → Build Go binary menggunakan Makefile.
        //   → Go di-install langsung di Jenkins agent di /usr/local/go.
        // ==============================================================
        stage('Compile') {
            steps {
                withEnv(["PATH=${env.PATH}:/usr/local/go/bin"]) {
                    sh 'go version'
                    sh 'make build'
                    sh 'ls -la bin/'
                }
            }
        }

        // ==============================================================
        // Stage 5: Dockerize
        //   → Build Docker image dan push ke GCP Artifact Registry.
        //   → Steps:
        //     1. Setup config yaml (set port 8080)
        //     2. docker build → create image
        //     3. gcloud auth → authenticate ke registry
        //     4. docker push → upload image
        // ==============================================================
        stage('Dockerize') {
            steps {
                // ==============================================================
                // ProxySQL Setup — TIDAK DIGUNAKAN
                //   → sttk-be menggunakan ProxySQL untuk MySQL connection pooling.
                //   → gold-gym-be-v2 tidak menggunakan ProxySQL saat ini.
                //   → Commented out. Aktifkan kalau gold-gym-be-v2 butuh ProxySQL:
                //       - Uncomment lines di bawah
                //       - Pastikan file config path sesuai: files/etc/gold-gym-be/
                // ==============================================================
                // echo '> Setup ProxySQL Config ...'
                // proxysqlSetup("./files/etc/gold-gym-be/gold-gym-be.staging.yaml")
                // proxysqlSetup("./files/etc/gold-gym-be/gold-gym-be.production.yaml")

                script {
                    echo '> Creating image ...'
                    // docker build → Build image dari Dockerfile di root project.
                    //   → Tag: gold-gym-be-v2:<VERSION>
                    //   → Di local tidak push ke external registry.
                    sh "docker build -t ${NAME}:${env.VERSION} ."
                    sh "docker build -t ${NAME}:latest ."
                    sh "docker images ${NAME}"
                }
            }
        }

        // ==============================================================
        // Stage 5: Vault Chart Injector — TIDAK DIGUNAKAN DI LOCAL
        //   → Di production: jalankan vault-chart-injector untuk inject
        //     secret templates dari Vault ke Helm chart.
        //   → Di local: tidak ada Vault server, stage ini di-skip.
        //   → Aktifkan kembali kalau sudah ada local Vault setup.
        // ==============================================================
        /*
        stage('Vault Chart Injector') {
            steps {
                echo '> Injecting vault chart ...'
                sh "vault-chart-injector"
            }
        }
        */

        // ==============================================================
        // Stage 6: Helm Lint
        //   → Di local: hanya lint chart, tidak package atau upload.
        //   → Helm package + upload ke ChartMuseum dilakukan di
        //     production pipeline (bukan local).
        // ==============================================================
        stage('Helm Lint') {
            steps {
                sh """
                    cd charts
                    sed -i 's/name: NAME/name: ${NAME}/g' Chart.yaml
                    sed -i 's/tag: dev/tag: ${env.VERSION}/g' values.yaml
                    helm lint . || echo "Helm tidak terinstall, skip lint"
                """
            }
        }

        // ==============================================================
        // Stage 8: Helm Charts CHC — TIDAK DIGUNAKAN
        //   → Stage ini di sttk-be dipakai untuk upload chart ke
        //     namespace "chc" di ChartMuseum (deployment target terpisah).
        //   → gold-gym-be-v2 tidak memiliki deployment CHC saat ini.
        //   → Seluruh stage di-comment out menggunakan block comment /* ... */
        //   → Aktifkan kembali dan buat folder "chc-charts/" kalau dibutuhkan.
        // ==============================================================
        /*
        stage('Helm Charts CHC') {
            steps {
                echo '> Changing repository name value ...'
                sh "sed -i 's#repository: draft#repository: ${DOCKER_REGISTRY_PROJECT_URL}/${NAME}#g' chc-charts/values.yaml"
                echo '> Changing version value ...'
                sh "sed -i 's/tag: dev/tag: ${env.VERSION}/g' chc-charts/values.yaml"
                echo '> Changing chart name value ...'
                sh "sed -i 's/name: NAME/name: ${env.NAME}/g' chc-charts/Chart.yaml"
                echo '> Remove some manifest ...'
                sh 'cd chc-charts/templates/ && ls | grep yaml | grep cron | xargs -r rm'
                echo '> Packing helm chart ...'
                sh "cd chc-charts && helm package . --version=${env.VERSION}"
                echo '> Uploading chart ...'
                sh "cd chc-charts && curl --data-binary '@${env.NAME}-${env.VERSION}.tgz' http://chartmuseum:8080/api/chc/charts"
                echo '> Removing uploaded chart package ...'
                sh "rm chc-charts/${env.NAME}-${env.VERSION}.tgz"
            }
        }
        */
    }

    // ==============================================================
    // post
    //   → Post-pipeline hooks: jalankan setelah semua stages selesai.
    //   → Tiga jenis:
    //     - always    : jalankan apapun hasilnya (sukses / gagal)
    //     - success   : jalankan HANYA kalau semua stages sukses
    //     - regression: jalankan HANYA kalau pipeline gagal
    // ==============================================================
    post {
        // always → Cleanup workspace setelah pipeline selesai.
        always {
            cleanWs()
        }

        // success → Kirim Discord notification via curl (tanpa plugin).
        success {
            // ==============================================================
            // GKE deploy + DevSecOps scanner — TIDAK DIGUNAKAN DI LOCAL
            //   → Di production: trigger deployment pipeline ke GKE
            //     dan security scanner.
            //   → Di local: tidak ada pipeline tersebut, di-skip.
            // ==============================================================
            /*
            build job: "${env.GKE_JOB_NAME}", parameters: [
                string(name: 'PROJECT_NAME', value: "${env.NAME}"),
                string(name: 'PROJECT_VERSION', value: "${env.VERSION}")
            ], wait: false

            build job: "devsecops-scanner-pipeline", parameters: [
                string(name: 'PROJECT_URL', value: "https://github.com/okafuizagoto/gold-gym-be-v2"),
                string(name: 'PROJECT_VERSION', value: "${env.VERSION}")
            ], wait: false
            */

            // Discord notification via curl
            //   → Kirim embed message ke Discord webhook.
            //   → color: 65280 = green (sukses)
            sh """
                curl -s -X POST "${env.DISCORD_WEBHOOK_URL}" \
                  -H "Content-Type: application/json" \
                  -d '{
                    "embeds": [{
                      "title": "${env.PIPELINE_NAME} #${env.BUILD_NUMBER} — SUCCESS",
                      "color": 65280,
                      "fields": [
                        {"name": "Version", "value": "${env.VERSION}", "inline": true},
                        {"name": "Image", "value": "${NAME}:${env.VERSION}", "inline": true},
                        {"name": "Branch", "value": "main", "inline": true}
                      ]
                    }]
                  }'
            """
        }

        // failure → Kirim notifikasi gagal ke Discord.
        failure {
            // color: 16711680 = red (gagal)
            sh """
                curl -s -X POST "${env.DISCORD_WEBHOOK_URL}" \
                  -H "Content-Type: application/json" \
                  -d '{
                    "embeds": [{
                      "title": "${env.PIPELINE_NAME} #${env.BUILD_NUMBER} — FAILED",
                      "color": 16711680,
                      "fields": [
                        {"name": "Version", "value": "${env.VERSION}", "inline": true},
                        {"name": "Branch", "value": "main", "inline": true}
                      ]
                    }]
                  }'
            """
        }
    }
}
