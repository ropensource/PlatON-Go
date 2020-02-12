package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SpecialFunctionsA;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

import java.text.SimpleDateFormat;
import java.util.Date;

/**
 *
 * @title 验证函数platon_block_number,platon_timestamp
 * @description:
 * @author: liweic
 * @create: 2020/02/10
 */
public class SpecialFunctionsATest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.SpecialFunctionsA验证链上区块相关函数",sourcePrefix = "wasm")
    public void SpecialfunctionsA() {

        try {
            prepare();
            SpecialFunctionsA specialfunctionsa = SpecialFunctionsA.deploy(web3j, transactionManager, provider).send();
            String contractAddress = specialfunctionsa.getContractAddress();
            String transactionHash = specialfunctionsa.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SpecialFunctionsATest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            Long blocknumber =specialfunctionsa.getBlockNumber().send();
            collector.logStepPass("getPlatONGas函数返回值:" + blocknumber);
            boolean result = "0".toString().equals(blocknumber.toString());
            collector.assertEqual(result, false);

            //bug
//            Long timestamp = specialfunctionsa.getTimestamp().send();
//            collector.logStepPass("block.timestamp函数返回值:" + timestamp);
//            SimpleDateFormat sdf=new SimpleDateFormat("yyyy-MM-dd");
//            String resultTime = sdf.format(new Date(Long.parseLong(String.valueOf(timestamp))));
//            System.out.print("时间：" + resultTime);
//            String today = sdf.format(new Date());
//            collector.assertEqual(resultTime,today);


        } catch (Exception e) {
            collector.logStepFail("SpecialFunctionsBTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
