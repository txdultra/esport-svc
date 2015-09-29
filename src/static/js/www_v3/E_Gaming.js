// JavaScript Document

 /*  writer:Yugc  */
 
 
 $(function(){
	 	if($.browser.msie && $.browser.version < 10){
		$('body').addClass('ltie10');
	}

	$('#dowebok').fullpage({
		verticalCentered: false,
		sectionsColor: ['#313131', '#ffffff', '#313131', '#FFF','#ffffff'],
		anchors: ['page1', 'page2', 'page3', 'page4', 'page5'],
		navigation:'true',
		navigationTooltips: ['首页','四大版块','赛事项目','随时随地','立即下载']

	});
	
	
    $(window).resize(function(){
        autoScrolling();
    });

    function autoScrolling(){   //宽度小于1024 出现滚动条，自适应
        var $ww = $(window).width();
        if($ww < 1024){
            $.fn.fullpage.setAutoScrolling(false);
        } else {
            $.fn.fullpage.setAutoScrolling(true);
        }
    }

    autoScrolling();
				
});

$(function(){
		$('.section').css('height',document.body.clientHeight);
		$(".box").hide();
		$('.wx_ma').hide();
		$('.qq_ma').hide();
		 
		 var maskHeight=$(document).height();
		 var maskWidth=$(document).width();
		
		$(".bug").click(function(){
			//添加遮罩层
			$('<div class="mask"></div>').appendTo($('body'));
			$('div.mask').css({
					'opacity':0.4,
					'background':'#000',
					'position':'absolute',
					'left':0,
					'top':0,
					'width':maskWidth,
					'height':maskHeight,
					'z-index':88
				});
			
			 $('.box').show();
			});
			
		$(".close").click(function(){
				$(".box").hide();
				$('.mask').remove();
			});
			
		$('.share1').hover(function(){
				$('.wx_ma').show();
			}
		,function(){
				$('.wx_ma').hide();
			});
			
		$('.share4').hover(function(){
				$('.qq_ma').show();
			},function(){
				$('.qq_ma').hide();
				});
});