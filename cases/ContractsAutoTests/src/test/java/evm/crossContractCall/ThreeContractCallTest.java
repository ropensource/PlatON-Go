package evm.crossContractCall;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.*;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title 0.5.13三个合约间跨合约调用 delegatecall只会修改第一个合约中的状态变量
 * @description:
 * @author: hudenian
 * @create: 2020/1/9
 */
public class ThreeContractCallTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ThreeContractCallTest-三个合约间跨合约调用者", sourcePrefix = "evm")
    public void threeContractCaller() {
        try {
            //第一个合约
            CallerOne callerOne = CallerOne.load("0x14d865935e5a5b7e1798bc8745625a33bdea3d24",web3j, transactionManager, provider);
            String callerContractAddress = callerOne.getContractAddress();


            //第二个合约
            CallerTwo callerTwo = CallerTwo.load("0x9f5839158882b3bffc2ee9198c5944c67f5093b8",web3j, transactionManager, provider);
            String callerTwoContractAddress = callerTwo.getContractAddress();

            //第三个合约
            CallerThree callerThree = CallerThree.load("0x533a923ffd5b7c83b48490a3ad939f75eb7fce62",web3j, transactionManager, provider);
            String callerThreeContractAddress = callerThree.getContractAddress();



            //查询第一个合约x值
            String callerOneX = callerOne.getCallerX().send().toString();
            collector.logStepPass("CallerOne 合约中X的值为："+callerOneX);

            //查询第二个合约x值
            String callerTwoX = callerTwo.getCalleeX().send().toString();
            collector.logStepPass("CallerTwo 合约中X的值为："+callerTwoX);

            //查询第三个合约x值
            String callerThreeX = callerThree.getCalleeThreeX().send().toString();
            collector.logStepPass("CallerThree 合约中X的值为："+callerThreeX);


            TransactionReceipt tx2 = callerOne.inc_delegatecall().send();
            collector.logStepPass("执行跨合约调用后，hash:" + tx2.getTransactionHash());

            //查询第一个合约x值
            callerOneX = callerOne.getCallerX().send().toString();
            collector.logStepPass("CallerOne 合约中X的值为："+callerOneX);
            collector.assertEqual("1",callerOneX);

            //查询第二个合约x值
            callerTwoX = callerTwo.getCalleeX().send().toString();
            collector.logStepPass("CallerTwo 合约中X的值为："+callerTwoX);
            collector.assertEqual("0",callerTwoX);

            //查询第三个合约x值
            callerThreeX = callerThree.getCalleeThreeX().send().toString();
            collector.logStepPass("CallerThree 合约中X的值为："+callerThreeX);
            collector.assertEqual("0",callerThreeX);

        } catch (Exception e) {
            collector.logStepFail("ThreeContractCallTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
