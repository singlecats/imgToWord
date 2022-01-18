$(function(){
  /*
   * For the sake keeping the code clean and the examples simple this file
   * contains only the plugin configuration & callbacks.
   * 
   * UI functions ui_* can be located in: demo-ui.js
   */
  $('#drag-and-drop-zone').dmUploader({ //
    url: '/upload',
    queue : false,
    extFilter : ['zip'],
    allowedTypes: "application/zip",
    maxFileSize: 30000000000, // 3 Megs
    onDragEnter: function(){
      // Happens when dragging something over the DnD area
      this.addClass('active');
    },
    onDragLeave: function(){
      // Happens when dragging something OUT of the DnD area
      this.removeClass('active');
    },
    onInit: function(){
      // Plugin is ready to use
      ui_add_log('初始化 :)', 'info');
    },
    onComplete: function(){
      // All files in the queue are processed (success or error)
      ui_add_log('准备完毕');
    },
    onNewFile: function(id, file){
      // When a new file is added using the file selector or the DnD area
      ui_add_log('新文件#' + id);
      ui_multi_add_file(id, file);
    },
    onBeforeUpload: function(id){
      // about tho start uploading a file
      ui_add_log('开始上传 #' + id);
      ui_multi_update_file_status(id, '上传中...', '上传中...');
      ui_multi_update_file_progress(id, 0, '', true);
    },
    onUploadCanceled: function(id) {
      // Happens when a file is directly canceled by the user.
      ui_multi_update_file_status(id, 'warning', '用户取消');
      ui_multi_update_file_progress(id, 0, 'warning', false);
    },
    onUploadProgress: function(id, percent){
      // Updating file progress
      ui_multi_update_file_progress(id, percent);
    },
    onUploadSuccess: function(id, data){
      // A file was successfully uploaded

      ui_add_log('上传结果 #' + id + ': ' + JSON.stringify(data));
      ui_add_log('文件 #' + id + ' 上传', '成功');
      if (data.Code != 200) {
        ui_multi_update_file_status(id, 'danger', '上传失败'+data.Message);
        ui_multi_update_file_progress(id, 100, 'danger', false);
      } else {
        ui_multi_update_file_status(id, 'success', '上传成功');
        ui_multi_update_file_progress(id, 100, 'success', false);
      }

    },
    onUploadError: function(id, xhr, status, message){
      ui_multi_update_file_status(id, 'danger', message);
      ui_multi_update_file_progress(id, 0, 'danger', false);  
    },
    onFallbackMode: function(){
      // When the browser doesn't support this plugin :(
      ui_add_log('Plugin cant be used here, running Fallback callback', 'danger');
    },
    onFileSizeError: function(file){
      ui_add_log('File \'' + file.name + '\' 文件过大', 'danger');
    }
  });
});