pipeline {
    agent {
        node {
            label "aws-nuclearport-jenkins"
        }
    }

    environment {
        BUILD_USER = ""
        robotid = "4520246e-8a86-4549-bf9b-15713250efd0"
        CRED_ID = "5d5f2b6e-fce1-4b60-b734-811b937918f7"
        REPO_URL = "git@gitlab.ftsview.com:ExternalProjects/TimeCapsuleStudio/game_slots_vsn.git"
    }

    post {
        success {
            script {
                currentBuild.description = "\n 构建成功!"
            }
            dingTalk(
                robot: robotid,
                type: 'ACTION_CARD',
                title: 'Deploy',
                text: [
                    "[${JOB_NAME}](${JOB_URL})\n",
                    '---',
                    "- 任务：[#${BUILD_ID}](${BUILD_URL})",
                    "- 状态：<font color=#52C41A>${currentBuild.currentResult}</font>",
                    "- 执行人：${BUILD_USER}",
                    "- 持续时间：${currentBuild.durationString}\n",
                ],
                btns: [
                    [
                        title    : '日志',
                        actionUrl: "${BUILD_URL}"
                    ],
                    [
                        title    : '控制台',
                        actionUrl: "https://jenkins.ftsview.com/"
                    ]
                ]
            )
        }

        failure {
            script {
                currentBuild.description = "\n 构建失败!"
            }
            dingTalk(
                robot: robotid,
                type: 'ACTION_CARD',
                title: 'Deploy',
                text: [
                    "[${JOB_NAME}](${JOB_URL})\n",
                    '---',
                    "- 任务：[#${BUILD_ID}](${BUILD_URL})",
                    "- 状态：<font color=#52C41A>${currentBuild.currentResult}</font>",
                    "- 执行人：${BUILD_USER}",
                    "- 持续时间：${currentBuild.durationString}\n",
                ],
                btns: [
                    [
                        title    : '日志',
                        actionUrl: "${BUILD_URL}"
                    ],
                    [
                        title    : '控制台',
                        actionUrl: "https://jenkins.ftsview.com/"
                    ]
                ]
            )
        }

        aborted {
            script {
                currentBuild.description = "\n 构建取消!"
            }
        }
    }

    parameters {
        choice(
                    name: "cluster",
                    choices: ['new'],
                    description: '选择目标集群'

            )
        gitParameter branchFilter: 'origin/(?!master$)(.*)', defaultValue: 'master', listSize: '10', name: 'branch_name', tagFilter: 'origin/master', sortMode: 'DESCENDING', type: 'PT_BRANCH'
    }



    stages {
        stage('获取依赖库代码') {
            steps {
                wrap([$class: 'BuildUser']) {
                   script {
                       BUILD_USER = "${env.BUILD_USER}"
                   }
                }
                script {
                    dir('./') {
                            checkout([$class                           : 'GitSCM',
                                      branches                         : [[name: "*/${params.branch_name}"]],
                                      doGenerateSubmoduleConfigurations: false,
                                      extensions                       : [],
                                      gitTool                          : 'Default',
                                      submoduleCfg                     : [],
                                      userRemoteConfigs                : [[credentialsId: '5d5f2b6e-fce1-4b60-b734-811b937918f7', url: env.REPO_URL]]
                            ])
                        }
                    }
                }
            }

        stage('构建镜像') {
                steps {
                    script {
                        //prod 生产环境
                        dir('./') {
                            sh "make build_stage env=${params.branch_name} ver=${BUILD_ID}"
                        }
                    }
                }
            }
            stage('推送构建镜像') {
                steps {
                    script {
                        dir('./') {
                           sh "make push_stage env=${params.branch_name} ver=${BUILD_ID}"
                        }
                    }
                }
            }
            stage('部署镜像') {
                steps {
                    script {
                        dir('./') {
                            if ("${params.cluster}" == "new") {
                                sh """
                                    /root/.kube/switch-role.sh
                                    # 使用 jq 提取并赋值到环境变量
                                    export AWS_ACCESS_KEY_ID=\$(jq -r '.Credentials.AccessKeyId' /tmp/assume_role_output.json)
                                    export AWS_SECRET_ACCESS_KEY=\$(jq -r '.Credentials.SecretAccessKey' /tmp/assume_role_output.json)
                                    export AWS_SESSION_TOKEN=\$(jq -r '.Credentials.SessionToken' /tmp/assume_role_output.json)

                                    make deploy_stage env=${params.branch_name} ver=${BUILD_ID} config=~/.kube/jack-EKS-1-31-prod

                                    unset AWS_ACCESS_KEY_ID
                                    unset AWS_SECRET_ACCESS_KEY
                                    unset AWS_SESSION_TOKEN
                                """

                            } else {
                                sh "make deploy_stage env=${params.branch_name} ver=${BUILD_ID} config=~/.kube/new-jackpland"
                            }
                        }
                    }
                }
            }
    }
}
