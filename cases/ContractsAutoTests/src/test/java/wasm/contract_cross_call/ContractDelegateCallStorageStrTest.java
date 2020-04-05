package wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDelegateCallStorageString;
import network.platon.contracts.wasm.ContractStorageString;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

public class ContractDelegateCallStorageStrTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_storage_str",sourcePrefix = "wasm")
    public void testDelegateCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `storge_str`, first
            ContractStorageString strc = ContractStorageString.load("0x46a6f6ef96765dd41b7a78ab016880581a4363cd",web3j, transactionManager, provider);


            String strcAddr = strc.getContractAddress();


            // deploy the delegate_call  contract second
            ContractDelegateCallStorageString delegateCall = ContractDelegateCallStorageString.load("0x902ecb5bdff08abbad05b59ed19b38a9ffb5505a",web3j, transactionManager, provider);

            String delegateCallAddr = delegateCall.getContractAddress();


            // check arr size 1st
            String strcStr = strc.get_string().send();
            collector.logStepPass("the msg count in arr of  storge_str contract:" + strcStr);
            collector.assertEqual(strcStr, "");

            String delegateCallStr = delegateCall.get_string().send();
            collector.logStepPass("the msg count in arr of cross_delegate_call_storage_str contract:" + delegateCallStr);
            collector.assertEqual(delegateCallStr, "");

            String msg = "Gavin";

            TransactionReceipt receipt = delegateCall.delegate_call_set_string(strcAddr, msg, Uint64.of(60000000l)).send();
            collector.logStepPass("cross_delegate_call_storage_str call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            strcStr = strc.get_string().send();
            collector.logStepPass("the msg count in arr of  storge_str contract:" + strcStr);
            collector.assertEqual(strcStr, "");

            delegateCallStr = delegateCall.get_string().send();
            collector.logStepPass("the msg count in arr of cross_delegate_call_storage_str contract:" + delegateCallStr);
            collector.assertEqual(delegateCallStr, msg);

        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_delegate_call_storage_str Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
