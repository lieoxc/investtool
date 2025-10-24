# InvesTool Frontend Build Script

# 构建前端项目
build_frontend() {
    echo "开始构建前端项目..."
    cd frontend
    
    # 安装依赖
    echo "安装依赖..."
    npm install
    
    # 构建项目
    echo "构建项目..."
    npm run build
    
    echo "前端构建完成！"
    cd ..
}

# 部署到 nginx
deploy_to_nginx() {
    echo "部署到 nginx..."
    
    # 复制构建文件到 nginx 目录
    sudo cp -r frontend/build/* /var/www/html/
    
    # 重启 nginx
    sudo systemctl restart nginx
    
    echo "部署完成！"
}

# 主函数
main() {
    case "$1" in
        "build")
            build_frontend
            ;;
        "deploy")
            build_frontend
            deploy_to_nginx
            ;;
        *)
            echo "用法: $0 {build|deploy}"
            echo "  build  - 构建前端项目"
            echo "  deploy - 构建并部署到 nginx"
            ;;
    esac
}

main "$@"
