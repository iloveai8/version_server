#!groovy

library 'fotoable-libs'
    def map = [:]

    //以下参数不需要修改
    //jenkins内置变量 job的名字
    def job_name = env.JOB_NAME.replaceAll("/","-")
    //构建分支，读取多分支构建的分支
    def BRANCH = env.BRANCH_NAME
    //以下参数需要研发人员修改
    //拉取代码库的地址
    map.put('REPO_URL',"git@gitlab.ftsview.com:ExternalProjects/TimeCapsuleStudio/game_slots_vsn.git")
    //代码的构建分支
    map.put('BRANCH', "${BRANCH}")
    //以下参数为多分支构建参数
    //测试环境
    if ("${BRANCH}" == "dev"){
        // 测试环境发版节点
        map.put('node','aws-nuclearport-jenkins')
        // 部署环境
        map.put('DEPENV','dev')
        map.put('cluster', "eks")
    //预发布环境
    } else if("${BRANCH}" == "pre"){
        // 预发布环境发版节点
        map.put('node','aws-nuclearport-jenkins')
        // 部署环境
        map.put('DEPENV','pre')
        map.put('cluster', "eks")
    //生产环境
    } else if ("${BRANCH}" == "pro"){
        // 生产环境发版节点
        map.put('node','aws-nuclearport-jenkins')
        // 部署环境
        map.put('DEPENV','pro')
        map.put('cluster', "eks")
    }
// 环境使用方法(dev为测试环境请使用k8s;stage为预发布使用ekst;master为生产环境请使用eks)
game_slots_vsn ("cluster",map)
