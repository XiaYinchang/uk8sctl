## uk8sctl

UCloud UK8S helper tool.

### 全局参数

全局参数
| 参数 | 解释 |
| ------------ | ----------------- |
| publickey | 用户在 UCloud 云平台生成的 PublicKey |
| privatekey | 用户在 UCloud 云平台生成的 PrivateKey |
| region | 指定资源所在地域，[点击查看地域列表](https://docs.ucloud.cn/api/summary/regionlist) |
| project-id | 指定资源所在项目，项目 Id 可在控制台首页查看到 |
| image-id | 指定镜像 ID 用于创建虚拟机，选填 |
| image-type | 指定镜像类型，可以是 centos 或者 ubuntu，默认为 centos， 可选填 |

### 支持的命令

当前仅支持 `create-base-uhost` 命令。

#### create-base-uhost

基于 UK8S 基础镜像创建一个虚拟机，用户可以在此虚拟机的基础上进行一些配置工作，之后通过虚拟机制作自定义镜像，镜像制作完成后可以将镜像 ImageId 传给 UK8S 相关 API 用于新建集群或者添加节点。

##### 参数列表

| 参数         | 是否必填 | 解释                                                             |
| ------------ | -------- | ---------------------------------------------------------------- |
| password     | 是       | 创建虚拟机时使用的密码 ，创建成功后用户使用该密码登录            |
| name         | 否       | 虚拟机名称，默认为 based-on-uk8s-image-for-creating-user-image   |
| cpu          | 否       | 虚拟机核数，默认为 2                                             |
| memory       | 否       | 虚拟机内存大小，默认为 4096                                      |
| charge-type  | 否       | 虚拟机收费类型，默认为 Month，按月收费                           |
| machine-type | 否       | 虚拟机类型，默认为 N ，代表普通机型，另有快杰机型等              |
| cpu-platform | 否       | 虚拟机处理器平台，默认为 Intel/Auto                              |
| disk-type    | 否       | 系统盘类型，默认为 CLOUD_SSD                                     |
| disk-size    | 否       | 系统盘大小，默认为 40，单位为 GB                                 |
| zone         | 否       | 在地域下指定可用区，若未指定，则在该地域所有可用区中随机选择一个 |
| vpcid        | 否       | 指定虚机所属 VPC ，若未指定，则在可用 VPC 中随机选择一个         |
| subnetid     | 否       | 指定虚机所属子网，若未指定，则在 VPC 所有子网中随机选择一个      |
| quantity     | 否       | 创建的虚机数量，默认为 1                                         |

#### 使用示例

执行命令：

```bash
./uk8sctl create-base-uhost --publickey your-public-key --privatekey you-private-key --region cn-bj2 --project-id org-test --image-type centos --password just-for-test
```

输出如下：

```bash
The following uhost instance is created based on uk8s image:
UhostId:   uhost-kjfssfss
Name:      based-on-uk8s-image-for-creating-user-image
ProjectId: org-test
Region:    cn-bj2
Zone:      cn-bj2-02
VPCId:     uvnet-test
SubnetId:  subnet-test
ImageId:   uimage-test
OsName:    UCloud_CentOS_7.6_UK8S_X.X
```

创建成功后，用户即可到 org-test 项目下北京二地域找到名称为 based-on-uk8s-image-for-creating-user-image 的虚拟机，并以此为基础制作自定义镜像。
